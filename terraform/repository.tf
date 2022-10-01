resource "aws_ecr_repository" "main" {
	name                 = "${var.name}"
	image_tag_mutability = "MUTABLE"
	tags                 = var.common_tags
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

resource "aws_ecr_repository_policy" "read_policy" {
  repository = aws_ecr_repository.main.name

  policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "AllowPull",
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "${aws_iam_role.ecs_task_execution_role.arn}",
					"arn:aws:iam::201843717406:root"
        ]
      },
      "Action": [
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:BatchCheckLayerAvailability",
				"ecr:*"
      ]
    }
  ]
}
EOF
}