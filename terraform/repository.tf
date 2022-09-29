resource "aws_ecr_repository" "main" {
	name                 = "${var.name}"
	image_tag_mutability = "MUTABLE"
	tags = var.common_tags
}

resource "aws_ecr_lifecycle_policy" "main" {
  repository = aws_ecr_repository.main.name
 
  policy = jsonencode({
    rules = [{
			rulePriority  = 1
			description   = "keep last 5 images"
			action        = {
				type        = "expire"
			}
			selection     = {
				tagStatus   = "any"
				countType   = "imageCountMoreThan"
				countNumber = 5
			}
    }]
  })
}