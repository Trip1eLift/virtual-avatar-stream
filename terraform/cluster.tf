resource "aws_ecs_cluster" "main" {
  name = "${var.name}-cluster-${var.environment}"
}

resource "aws_ecs_task_definition" "main" {
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  container_definitions = jsonencode([{
		name        = "${var.name}-container-${var.environment}"
		image       = "${var.container_image}:latest"
		essential   = true
		environment = var.container_environment
		portMappings = [{
			protocol      = "tcp"
			containerPort = var.container_port
			hostPort      = var.container_port
		}]
	}])
}

resource "aws_ecs_service" "main" {
	name                               = "${var.name}-service-${var.environment}"
	cluster                            = aws_ecs_cluster.main.id
	task_definition                    = aws_ecs_task_definition.main.arn
	desired_count                      = 2
	deployment_minimum_healthy_percent = 50
	deployment_maximum_percent         = 200
	launch_type                        = "FARGATE"
	scheduling_strategy                = "REPLICA"
	
	network_configuration {
		security_groups  = var.ecs_service_security_groups
		subnets          = var.subnets.*.id
		assign_public_ip = true
	}
	
	# load_balancer {
	# 	target_group_arn = var.aws_alb_target_group_arn
	# 	container_name   = "${var.name}-container-${var.environment}"
	# 	container_port   = var.container_port
	# }
	
	lifecycle {
		ignore_changes = [task_definition, desired_count]
	}
}