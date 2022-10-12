resource "aws_iam_role" "ecs_task_role" {
	name = "${var.name}-ecsTaskRole"
	tags = var.common_tags
	
	assume_role_policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Action": "sts:AssumeRole",
			"Principal": {
				"Service": "ecs-tasks.amazonaws.com"
			},
			"Effect": "Allow",
			"Sid": ""
		}
	]
}
	EOF
}

resource "aws_iam_role" "ecs_task_execution_role" {
	name = "${var.name}-ecsTaskExecutionRole"
	tags = var.common_tags
	
	assume_role_policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Action": "sts:AssumeRole",
			"Principal": {
				"Service": "ecs-tasks.amazonaws.com"
			},
			"Effect": "Allow",
			"Sid": ""
		}
	]
}
EOF
}
 
resource "aws_iam_role_policy_attachment" "ecs-task-execution-role-policy-attachment" {
	role       = aws_iam_role.ecs_task_execution_role.name
	policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# aws_iam_role.ecs_task_role.name does not have any policy attached to it.