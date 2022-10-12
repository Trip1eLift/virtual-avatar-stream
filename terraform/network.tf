resource "aws_vpc" "main" {
	cidr_block = var.vpc_cidr
	tags       = var.common_tags
}

resource "aws_internet_gateway" "main" {
	vpc_id = aws_vpc.main.id
	tags   = var.common_tags
}

resource "aws_subnet" "public" {
	vpc_id                  = aws_vpc.main.id
	cidr_block              = element(var.subnet_cidrs, count.index)
	availability_zone       = element(var.availability_zones, count.index)
	count                   = length(var.subnet_cidrs)
	map_public_ip_on_launch = true
	tags                    = var.common_tags
}

resource "aws_route_table" "public" {
	vpc_id = aws_vpc.main.id
	tags   = var.common_tags
}

resource "aws_route" "public" {
	route_table_id         = aws_route_table.public.id
	destination_cidr_block = "0.0.0.0/0"
	gateway_id             = aws_internet_gateway.main.id
}

resource "aws_route_table_association" "public" {
	count          = length(var.availability_zones)
	subnet_id      = element(aws_subnet.public.*.id, count.index)
	route_table_id = aws_route_table.public.id
}

resource "aws_security_group" "ecs_tasks" {
	name   = "${var.name}-${var.environment}-sg"
	vpc_id = aws_vpc.main.id
	
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