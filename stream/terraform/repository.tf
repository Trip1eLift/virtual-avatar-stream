resource "aws_ecr_repository" "main" {
	name                 = "${var.name}"
	image_tag_mutability = "MUTABLE"
	tags                 = var.common_tags
}

resource "aws_ecr_lifecycle_policy" "main" {
	repository = aws_ecr_repository.main.name
	
	policy = jsonencode({
		rules = [{
				rulePriority    = 1
				description     = "keep last 5 images"
				action          = {
					type        = "expire"
				}
				selection       = {
					tagStatus   = "any"
					countType   = "imageCountMoreThan"
					countNumber = 5
				}
		}]
	})
}

resource "null_resource" "docker_build_push" {
	depends_on = [
		aws_ecr_repository.main
	]

	triggers = {
		always_run = "${timestamp()}"
	}

	# These scripts require aws cli and a running docker
	provisioner "local-exec" {
		command = "aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin ${var.aws_account_id}.dkr.ecr.us-east-1.amazonaws.com"
	}
	provisioner "local-exec" {
		command = "docker build -t ${var.name} ../go-server"
	}
	provisioner "local-exec" {
		command = "docker tag ${var.name}:latest ${var.aws_account_id}.dkr.ecr.us-east-1.amazonaws.com/${var.name}:latest"
	}
	provisioner "local-exec" {
		command = "docker push ${var.aws_account_id}.dkr.ecr.us-east-1.amazonaws.com/${var.name}:latest"
	}
}