# route53-myip

route53-myip is a Go script that creates a Route53 A record pointing to my current public IP address. It is meant to be run periodically as a cron job in my Raspberry pi.

## The script

The Go script works by fetching my current public IP address using the [ipify API](https://ipify.org). It will compare this IP with the IP stored in a local `ip.txt` file that is used as a cache to avoid invoking Route53 even if the IP hasn't change. If the IP indeed has changed, then it will `UPSERT` a Route53 `A` record with the value of the new IP.

### Environment variables

The script is highly customizable by using the following environment variables

| Environment variable name | Description                                               | Default value              |
| ------------------------- | --------------------------------------------------------- | -------------------------- |
| ROUTE53_MYIP_DEBUG        | Enable debug logs                                         | `false`                    |
| ROUTE53_MYIP_LOG_PATH     | Path to the log file                                      | `"stdout"`                 |
| ROUTE53_MYIP_REGION       | AWS region                                                | `"eu-south-2"`             |
| ROUTE53_MYIP_NAME         | The name of the A record                                  | `"home.ignaciovaquero.es"` |
| ROUTE53_MYIP_FILE_PATH    | The path to the local `ip.txt` file to be used as a cache | `"./ip.txt"`               |
| ROUTE53_MYIP_IPIFY_URL    | The URL to the ipify API                                  | `"https://api.ipify.org"`  |
| ROUTE53_MYIP_TTL          | The TTL of the A record in seconds                        | `1800`                     |


## Installation in Raspberry Pi

The script comes with some Terraform and Ansible definitions to make installation and uninstallation easy and completely automated. At the time this was created, we tested it with the following versions:
- Terraform: 1.3.9
- Ansible: 7.2.0

Both Terraform and Ansible code are idempotent.

### Terraform

Terraform uses the [aws](https://registry.terraform.io/providers/hashicorp/aws/4.55.0) provider to create the IAM user with the required inline policy to create the Route53 A record. The policy looks like this:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "route53:ChangeResourceRecordSets",
      "Condition": {
        "StringEquals": {
          "route53:ChangeResourceRecordSetsActions": "UPSERT",
          "route53:ChangeResourceRecordSetsNormalizedRecordNames": "home.ignaciovaquero.es",
          "route53:ChangeResourceRecordSetsRecordTypes": "A"
        }
      },
      "Effect": "Allow",
      "Resource": "arn:aws:route53:::hostedzone/Z076159532BK132ZKZOI9",
      "Sid": "VisualEditor0"
    },
    {
      "Action": "route53:ListHostedZonesByName",
      "Effect": "Allow",
      "Resource": "*",
      "Sid": "VisualEditor1"
    }
  ]
}
```

The permissions are quite tight, so that only the specific `home.ignaciovaquero.es` record name, of type `A` can be created in the proper hosted zone.

The Terraform code has a final [local_sensitive_file](https://registry.terraform.io/providers/hashicorp/local/2.3.0/docs/resources/sensitive_file) resource that creates the AWS [credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html#cli-configure-files-where) file and places it under the `ansible/roles/myip/files/credentials` path, so that Ansible can copy it to the Raspberry afterwards.

### Ansible

Ansible takes care of all the installation that happens in the Raspberry pi. Specifically, it:
- Creates the `myip` user and group
- Creates the proper folders:
  - /opt/myip/bin
  - /var/opt/myip
  - /var/log/myip
  - /home/myip/.aws
- Copies the binary file to the proper location (`/opt/myip/bin/myip`)
- Copies the aws credentials file to the proper location (`/home/myip/.aws/credentials`)
- Creates the cron file at `/etc/cron.d/myip`. **Currently the cron periodicity is hard coded and it will run every 2 hours**.

Ansible also takes care of the uninstallation steps. The steps involved in the uninstallation are:
- Deleting the cron file
- Deleting the myip user and group
- Deleting all the files and directories:
  - /opt/myip
  - /var/opt/myip
  - /home/myip
  - /var/log/myip

There is a global `myip_action` variable that handles installation and uninstallation. This variable can only take two values: `install` and `uninstall`, and it defaults to `install`. However, all this should be transparent to you if you use the `Makefile` (see next section).

### How all this ties together: The Makefile

Since there are multiple pieces involved in the installation/uninstallation of the script, we created a `Makefile` to properly handle everything.

#### Install

In order to call the installation steps, simply run:
```bash
make install
```

This will first build the Go code from source (and handle the proper `GOOS` and `GOARCH` environment variables to target the Raspberry pi architecture and OS), then it will invoke the Terraform code to create the IAM user, and it will lastly invoke the Ansible code. It will finally delete the `myip` binary generating during build and the `credentials` file generated by Terraform from the `ansible/roles/myip/files` directory.

#### Uninstall

In order to call the uninstallation steps, simply run:
```bash
make uninstall
```

This will first call the Ansible code, passing the `-e "myip_action=uninstall"` flag to actually uninstall the script from the Raspberry. Finally, it will call the Terraform code to destroy all the Terraform resources created.
