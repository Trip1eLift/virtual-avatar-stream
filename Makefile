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

deploy: deploy-init
	cd terraform; terraform apply
# after deployment, manually stop tasks and the new tasks will use the new image

destroy:
	cd terraform; terraform apply -destroy