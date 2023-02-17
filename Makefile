.PHONY: clean build tf_init tf_plan tf_apply tf_destroy install uninstall

clean:
	rm -f ansible/roles/myip/files/myip
	rm -f ip.txt

build:
	@go mod tidy
	GOOS=linux GOARCH=arm go build -o ansible/roles/myip/files/myip

tf_init:
	@terraform -chdir=terraform init

tf_plan: tf_init
	@terraform -chdir=terraform plan

tf_apply: tf_init
	@terraform -chdir=terraform apply -auto-approve

tf_destroy: tf_init
	@terraform -chdir=terraform destroy -auto-approve

install: build tf_init tf_apply
	ansible-playbook -i ansible/hosts ansible/main.yml
	@rm -f ansible/roles/myip/files/myip
	@rm -f ansible/roles/myip/files/credentials

uninstall:
	ansible-playbook -i ansible/hosts -e "myip_action=uninstall" ansible/main.yml
	@terraform -chdir=terraform destroy -auto-approve

