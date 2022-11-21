hosted_zone_id = "Z02316541LCC6J0N3IECG"
domain_name = "virtualavatar-stream.trip1elift.com"

name = "virtual-avatar-stream"
environment = "dev"

vpc_cidr = "10.0.0.0/16"
public_subnet_cidrs =  ["10.0.0.0/24", "10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.8.0/24", "10.0.9.0/24", "10.0.10.0/24"]

availability_zones =   ["us-east-1a", "us-east-1b", "us-east-1c"]

container_port = 5000

common_tags = {
	Project    = "virtualavatar"
	Owner      = "Trip1eLift"
	Repository = "https://github.com/Trip1eLift/Virtual-Avatar-Streaming-Backend"
	Management = "Managed by Terraform" 
}

aws_account_id = "201843717406"

database_settings = {
	DB_HOST = "postgres_service"
	DB_USER = "postgres_user"
	DB_NAME = "virtualavatar"
	DB_PORT = 5432
}

frontend_origin = "http://localhost:3000"
