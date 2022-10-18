variable "name" {
	type = string
  	description = "main name of the application"
}

variable "environment" {
	type = string
  	description = "environment of the deployment"
}

variable "vpc_cidr" {
	type = string
	description = "cidr block of vpc"
}

variable "subnet_cidrs" {
	type = list
	description = "cidr block of subnet (list length shoud be the same as availability_zones)"
}

variable "availability_zones" {
	type = list
	description = "list of availability zones for subnet (list length shoud be the same as subnet_cidrs)"
}

variable "container_port" {
	type = number
	description = "exposing port number of the docker container"
}

variable "cloudwatch_group" {
	type = string
	description = "cloudwatch logging group name of ecs task"
}

variable "common_tags" {
	type = map
	description = "common tags you want applied to all components"
}

variable "aws_account_id" {
	type = string
	description = "aws account id for deployment"
}