.PHONY: clean build tf_init tf_plan tf_apply tf_destroy create_k8s_secret delete_k8s_secret install uninstall

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

create_k8s_secret: tf_apply
	kubectl apply -k deploy/production/secret

delete_k8s_secret:
	kubectl delete -k deploy/production/secret

install: build tf_apply
	ansible-playbook -i ansible/hosts ansible/main.yml $(add_hosts)
	kubectl apply -k deploy/production/secret
	kubectl apply -k deploy/production
	@rm -f ansible/roles/myip/files/myip
	@rm -f ansible/roles/myip/files/credentials
	@rm -f deploy/production/secret/credentials

uninstall:
	ansible-playbook -i ansible/hosts -e "myip_action=uninstall" ansible/main.yml $(add_hosts)
	k delete -k deploy/production
	k delete -k deploy/production/secret
	@terraform -chdir=terraform destroy -auto-approve
