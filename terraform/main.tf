terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.55.0"
    }
  }

  backend "s3" {
    bucket = "igvaquero-terraform-state"
    key    = "route53-myip/state"
    region = "eu-west-1"
  }
}

provider "aws" {
  region = var.region
}

resource "aws_iam_user" "myip" {
  name          = "myip"
  path          = "/raspberry/"
  force_destroy = true

  tags = {
    device = "raspberry"
  }
}

resource "aws_iam_access_key" "myip" {
  user = aws_iam_user.myip.name
}

resource "aws_iam_user_policy" "route53_myip" {
  name = "MyIP_Route53_permissions"
  user = aws_iam_user.myip.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "VisualEditor0"
        Effect   = "Allow"
        Action   = "route53:ChangeResourceRecordSets"
        Resource = "arn:aws:route53:::hostedzone/Z076159532BK132ZKZOI9"
        Condition = {
          StringEquals = {
            "route53:ChangeResourceRecordSetsRecordTypes"           = var.record_conditions.type
            "route53:ChangeResourceRecordSetsNormalizedRecordNames" = var.record_conditions.name
            "route53:ChangeResourceRecordSetsActions"               = var.record_conditions.action
          }
        }
      },
      {
        Sid      = "VisualEditor1"
        Effect   = "Allow"
        Action   = "route53:ListHostedZonesByName"
        Resource = "*"
      }
    ]
    }
  )
}

resource "local_sensitive_file" "aws_credentials" {
  filename        = "${path.root}/../ansible/roles/myip/files/credentials"
  content         = <<EOF
[default]
aws_access_key_id=${aws_iam_access_key.myip.id}
aws_secret_access_key=${aws_iam_access_key.myip.secret}
EOF
  file_permission = 0444
}
