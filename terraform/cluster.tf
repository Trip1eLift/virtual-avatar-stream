resource "aws_ecs_cluster" "main" {
  name = "${var.name}-${var.environment}-cluster"
	tags = var.common_tags
}

resource "aws_ecs_task_definition" "main" {
	family                   = "${var.name}-${var.environment}-family"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  container_definitions = jsonencode([{
		name        = "${var.name}-${var.environment}-container"
		image       = "${var.container_image}:latest"
		essential   = true
		environment = [{"name": "environment", "value": "${var.container_environment}"}]
		portMappings = [{
			protocol      = "tcp"
			containerPort = var.container_port
			hostPort      = var.container_port
		}]
		# TODO: Need to handle event for health check
		# healthCheck = {
    #   retries = 10
    #   command = [ "CMD-SHELL", "curl -f http://localhost:8081/actuator/liveness || exit 1" ]
    #   timeout: 5
    #   interval: 10
    #   startPeriod: var.health_start_period
    # }
		logConfiguration = {
      logDriver = "awslogs"
      options   = {
        awslogs-group         = var.cloudwatch_group # TODO: need to create cloudwatch group somewhere and reference it here
        awslogs-region        = "us-east-1"
        awslogs-stream-prefix = "ecs"
      }
    }
	}])
	tags = var.common_tags
}

resource "aws_ecs_service" "main" {
	name                               = "${var.name}-${var.environment}-service"
	cluster                            = aws_ecs_cluster.main.id
	task_definition                    = aws_ecs_task_definition.main.arn
	desired_count                      = 2
	deployment_minimum_healthy_percent = 50
	deployment_maximum_percent         = 200
	launch_type                        = "FARGATE"
	scheduling_strategy                = "REPLICA"

	network_configuration {
		security_groups  = [ aws_security_group.ecs_tasks.id ]
  	subnets          = flatten(aws_subnet.public.*.id)
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
	tags = var.common_tags
}