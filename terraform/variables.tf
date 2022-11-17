variable "hosted_zone_id" {
  description = "Id of the pre-defined hosted zone."
}

variable "domain_name" {
	type = string
	description = "domain name (including sub-domain)"
}

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

variable "public_subnet_cidrs" {
	type = list
	description = "public cidr block of subnet (list length shoud be the same as availability_zones)"
}

variable "private_subnet_cidrs" {
	type = list
	description = "private cidr block of subnet (list length shoud be the same as availability_zones)"
}

variable "availability_zones" {
	type = list
	description = "list of availability zones for subnet (list length shoud be the same as subnet_cidrs)"
}

variable "container_port" {
	type = number
	description = "exposing port number of the docker container"
}

variable "common_tags" {
	type = map
	description = "common tags you want applied to all components"
}

variable "aws_account_id" {
	type = string
	description = "aws account id for deployment"
}