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
		image       = "${aws_ecr_repository.main.repository_url}:latest"
		essential   = true
		environment = [{"name": "environment", "value": "${var.environment}"}]
		portMappings = [{
			protocol      = "tcp"
			containerPort = var.container_port
			hostPort      = var.container_port
		}]
		# Healthcheck is written in Dockerfile.
		# healthCheck = {
		# 	command     = [ "CMD-SHELL", "curl -sf http://localhost:5001/health || exit 1" ]
		# 	retries     = 3
		# 	timeout     = 3
		# 	interval    = 5
		# 	startPeriod = 5
    	# }
		logConfiguration = {
			logDriver = "awslogs"
			options   = {
				awslogs-group         = aws_cloudwatch_log_group.main.name
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
	platform_version                   = "1.4.0"
	force_new_deployment               = true

	network_configuration {
		security_groups  = [ aws_security_group.ecs_tasks.id ]
  		subnets          = flatten(aws_subnet.public.*.id)
		assign_public_ip = true
	}
	
	# disable alb because it costs $30+ a month
	# load_balancer {
	# 	target_group_arn = aws_alb_target_group.main.arn
	# 	container_name   = "${var.name}-${var.environment}-container"
	# 	container_port   = var.container_port
	# }
	
	# desired_count is dynamic based on the scaling policies
	# force update desired_count to a higher number can achieve blue/green deployment
	lifecycle {
		ignore_changes = [desired_count]
	}
	tags = var.common_tags
}

resource "aws_cloudwatch_log_group" "main" {
	name = "${var.name}-log"
	tags = var.common_tags
}