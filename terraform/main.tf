terraform {
	required_version = "~> 1.2.9"

	required_providers {
		aws = {
			source = "hashicorp/aws"
			version = "~> 3.0"
		}
	}

	backend "s3" {
		bucket = "trip1elift-terraform"
		key = "virtualavatar-backend/terraform.tfstate"
		region = "us-east-1"
	}
}

provider "aws" {
    region = "us-east-1"
}