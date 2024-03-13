.PHONY: clean build tf_init tf_plan tf_apply tf_destroy install uninstall

HOSTS ?=
ifdef HOSTS
add_hosts = --limit $(HOSTS)
endif
VERSION := 1.0.0


clean:
	rm -f ansible/roles/myip/files/myip
	rm -f ansible/roles/myip/files/credentials
	rm -f deploy/production/credentials
	rm -f ip.txt

build:
	@go mod tidy
	GOOS=linux GOARCH=arm go build -o ansible/roles/myip/files/myip
	docker build --tag ivaquero/myip:$(VERSION) --tag ivaquero/myip:latest .
	docker push ivaquero/myip:$(VERSION)
	docker push ivaquero/myip:latest

tf_init:
	@terraform -chdir=terraform init

tf_plan: tf_init
	@terraform -chdir=terraform plan

tf_apply: tf_init
	@terraform -chdir=terraform apply -auto-approve

tf_destroy: tf_init
	@terraform -chdir=terraform destroy -auto-approve

install: build tf_init tf_apply
	ansible-playbook -i ansible/hosts ansible/main.yml $(add_hosts)
	k apply -k deploy/production
	@rm -f ansible/roles/myip/files/myip
	@rm -f ansible/roles/myip/files/credentials
	@rm -f deploy/production/credentials

uninstall:
	ansible-playbook -i ansible/hosts -e "myip_action=uninstall" ansible/main.yml $(add_hosts)
	k delete -k deploy/production
	@terraform -chdir=terraform destroy -auto-approve
