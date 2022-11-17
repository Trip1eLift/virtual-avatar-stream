hosted_zone_id = "Z02316541LCC6J0N3IECG"
domain_name = "virtualavatar-stream.trip1elift.com"

name = "virtual-avatar-stream"
environment = "dev"

vpc_cidr = "10.0.0.0/16"
public_subnet_cidrs = ["10.0.0.0/24", "10.0.1.0/24"]
private_subnet_cidrs = ["10.0.8.0/24", "10.0.9.0/24"]

availability_zones = ["us-east-1a", "us-east-1b"]

container_port = 5000

common_tags = {
	Project    = "virtualavatar"
	Owner      = "Trip1eLift"
	Repository = "https://github.com/Trip1eLift/Virtual-Avatar-Streaming-Backend"
	Management = "Managed by Terraform" 
}

aws_account_id = "201843717406"