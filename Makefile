start:
	docker-compose up --build

start-bg:
	docker-compose up --build --detach

down:
	docker-compose down -v

cleanse:
	docker system prune -a && docker volume prune

deploy-init:
	cd terraform; terraform init

deploy:
	cd terraform; terraform apply
# After deployment, manually stop old tasks. The new tasks will be deployed already.

cheap:
	cd terraform-cheap-deployment; terraform apply

destroy:
	cd terraform; terraform apply -destroy