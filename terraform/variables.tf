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

variable "container_image" {
	type = string
  description = "image name of the container"
}

variable "container_environment" {
	type = string
  description = "image environment of the container"
}

variable "container_port" {
	type = number
  description = "exposing port number of the docker container"
}

variable "common_tags" {
	type = map
	description = "Common tags you want applied to all components."
}