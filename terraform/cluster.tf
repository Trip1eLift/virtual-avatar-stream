resource "aws_ecs_cluster" "main" {
	name = "${var.name}-${var.environment}-cluster"
	tags = var.common_tags
}

resource "aws_ecs_task_definition" "main" {
	family                   = "${var.name}-${var.environment}"
	network_mode             = "awsvpc"
	requires_compatibilities = ["FARGATE"]
	cpu                      = 512
	memory                   = 1024
	execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
	task_role_arn            = aws_iam_role.ecs_task_role.arn
	container_definitions = jsonencode([{
		name        = "${var.name}-${var.environment}-container"
		image       = "${aws_ecr_repository.main.repository_url}:latest"
		essential   = true
		environment = [
			{"name": "environment",      "value": "${var.environment}"},
			{"name": "PORT",             "value": "${tostring(var.container_port)}"},
			{"name": "DB_HOST",          "value": "${aws_rds_cluster.main.endpoint}"},
			{"name": "DB_USER",          "value": "${aws_rds_cluster.main.master_username}"},
			{"name": "DB_NAME",          "value": "${aws_rds_cluster.main.database_name}"},
			{"name": "DB_PORT",          "value": "${tostring(aws_rds_cluster.main.port)}"},
			{"name": "DB_PASS",          "value": "postgres_password"},       # TODO: use AWS secret manager later
			{"name": "DB_RETRY_BACKOFF", "value": "${var.database_settings.DB_RETRY_BACKOFF}"},
			{"name": "ORIGIN_LOCAL",     "value": "${var.frontend_origin_local}"},
			{"name": "ORIGIN_REMOTE",    "value": "${var.frontend_origin_remote}"},
			{"name": "AISLE_KEY",        "value": "passcode"},                # TODO: use AWS secret manager later
			{"name": "TIME_STAMP",       "value": "${timestamp()}"},          # Forces update on image version which auto trigger deployment on new task
		]
		portMappings = [{
			protocol      = "tcp"
			containerPort = var.container_port
			hostPort      = var.container_port
		}]
		# Container health checks
		healthCheck = {
			command     = [ "CMD-SHELL", "curl -sf http://localhost:${var.container_port}/health || exit 1" ]
			retries     = 3
			timeout     = 3
			interval    = 5
			startPeriod = 5
		}
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
	desired_count                      = 1
	# deployment_minimum_healthy_percent = 100
	# deployment_maximum_percent         = 400
	launch_type                        = "FARGATE"
	scheduling_strategy                = "REPLICA"
	platform_version                   = "1.4.0"
	force_new_deployment               = true

	network_configuration {
		security_groups  = [ aws_security_group.ecs_service.id ]
		subnets          = flatten(aws_subnet.private.*.id)
		assign_public_ip = false
	}
	
	load_balancer {
		target_group_arn = aws_alb_target_group.main.arn
		container_name   = "${var.name}-${var.environment}-container"
		container_port   = var.container_port
	}
	
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