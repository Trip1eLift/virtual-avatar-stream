resource "aws_security_group" "alb" {
  name        = "${var.name}-${var.environment}-alb-sg"
  description = "For alb"
  vpc_id      = aws_vpc.main.id
 
  ingress {
   protocol         = "tcp"
   from_port        = 80
   to_port          = 80
   cidr_blocks      = ["0.0.0.0/0"]
   ipv6_cidr_blocks = ["::/0"]
  }
 
  ingress {
   protocol         = "tcp"
   from_port        = 443
   to_port          = 443
   cidr_blocks      = ["0.0.0.0/0"]
   ipv6_cidr_blocks = ["::/0"]
  }
 
  egress {
   protocol         = "-1"
   from_port        = 0
   to_port          = 0
   cidr_blocks      = ["0.0.0.0/0"]
   ipv6_cidr_blocks = ["::/0"]
  }

	tags = var.common_tags
}

resource "aws_security_group" "ecs_service" {
	name        = "${var.name}-${var.environment}-ecs-sg"
  description = "For ecs service"
	vpc_id      = aws_vpc.main.id
	
	ingress {
		protocol         = "tcp"
		from_port        = var.container_port
		to_port          = var.container_port
		cidr_blocks      = ["0.0.0.0/0"]
		ipv6_cidr_blocks = ["::/0"]
	}

	egress {
		# egress all
		protocol         = "-1"
		from_port        = 0
		to_port          = 0
		cidr_blocks      = ["0.0.0.0/0"]
		ipv6_cidr_blocks = ["::/0"]
	}

	tags = var.common_tags
}

# TODO: figure out where to attach this
resource "aws_security_group" "aurora" {
  name        = "${var.name}-${var.environment}-aurora-sg"
  description = "For aurora"
  vpc_id      = aws_vpc.main.id

  ingress {
		protocol         = "tcp"
		from_port        = var.database_settings.DB_PORT
		to_port          = var.database_settings.DB_PORT
		cidr_blocks      = ["0.0.0.0/0"]
		ipv6_cidr_blocks = ["::/0"]
	}

  tags = var.common_tags
}