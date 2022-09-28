variable "name" {
	type = string
  description = "main name of the application"
}

variable "environment" {
	type = string
  description = "environment of the deployment"
}

variable "cidr" {
  type = string
  description = "cidr block of vpc"
}

variable "public_subnets" {
	type = string
  description = ""
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

variable "ecs_service_security_groups" {
	type = string
  description = ""
}

variable "subnets" {
	type = string
  description = ""
}

variable "common_tags" {
	type = map
	description = "Common tags you want applied to all components."
}