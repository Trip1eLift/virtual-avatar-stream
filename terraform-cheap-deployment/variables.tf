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
	description = "public cidr block of subnet; ALB requires at least 2 subnets (list length shoud be the same as public_availability_zones)"
}

# variable "private_subnet_cidrs" {
# 	type = list
# 	description = "private cidr block of subnet; RDS requires at least 3 subnets (list length shoud be the same as private_availability_zones)"
# }

variable "public_availability_zones" {
	type = list
	description = "list of availability zones for public subnets (list length shoud be the same as public_subnet_cidrs)"
}

# variable "private_availability_zones" {
# 	type = list
# 	description = "list of availability zones for private subnets (list length shoud be the same as private_subnet_cidrs)"
# }

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

variable "database_settings" {
	type = map
	description = "settings of database"
}

variable "frontend_origin_local" {
	type = string
	description = "allowed origin of frontend for websocket connections"
}

variable "frontend_origin_remote" {
	type = string
	description = "allowed origin of frontend for websocket connections"
}