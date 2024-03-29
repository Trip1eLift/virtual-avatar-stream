hosted_zone_id = "Z02316541LCC6J0N3IECG"
domain_name = "virtualavatar-stream.trip1elift.com"

name = "virtual-avatar-stream"
environment = "cheap"

vpc_cidr = "10.0.0.0/16"
public_subnet_cidrs =        ["10.0.0.0/24", "10.0.1.0/24"]
public_availability_zones =  ["us-east-1a",  "us-east-1b" ]

container_port = 5000

common_tags = {
	Project    = "virtualavatar"
	Owner      = "Trip1eLift"
	Repository = "https://github.com/Trip1eLift/Virtual-Avatar-Streaming-Backend"
	Management = "Managed by Terraform"
	Version    = "Cheap deployment"
}

aws_account_id = "201843717406"

database_settings = {
	DB_HOST          = "postgres_service"
	DB_USER          = "postgres_user"
	DB_NAME          = "virtualavatar"
	DB_PORT          = 5432
	DB_RETRY_BACKOFF = 60
}

frontend_origin_local = "http://localhost:3000"
frontend_origin_remote = "https://virtualavatar.trip1elift.com"
# TODO: switch origin once in prod
