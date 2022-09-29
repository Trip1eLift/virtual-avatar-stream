name = "virtual-avatar-stream"
environment = "dev"
vpc_cidr = "10.0.0.0/16"
subnet_cidrs = ["10.0.0.0/24", "10.0.1.0/24"]
availability_zones = ["us-east-1a", "us-east-1b"]
container_image = "virtual-avatar-stream"
container_environment = "dev"
container_port = 5001
cloudwatch_group = "virtual-avatar-stream-log"

common_tags = {
  Project    = "virtualavatar"
  Owner      = "Trip1eLift"
  Repository = "https://github.com/Trip1eLift/Virtual-Avatar-Streaming-Backend"
  Management = "Managed by Terraform" 
}