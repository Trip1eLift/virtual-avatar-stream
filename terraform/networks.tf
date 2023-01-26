resource "aws_vpc" "main" {
	cidr_block = var.vpc_cidr
	tags       = var.common_tags
}

resource "aws_internet_gateway" "main" {
	vpc_id = aws_vpc.main.id
	tags   = var.common_tags
}

# Public subnets

resource "aws_subnet" "public" {
	vpc_id                  = aws_vpc.main.id
	cidr_block              = element(var.public_subnet_cidrs, count.index)
	availability_zone       = element(var.public_availability_zones, count.index)
	count                   = length(var.public_subnet_cidrs)
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
	count          = length(var.public_subnet_cidrs)
	subnet_id      = element(aws_subnet.public.*.id, count.index)
	route_table_id = aws_route_table.public.id
}

# Private subnets

resource "aws_subnet" "private" {
	vpc_id                  = aws_vpc.main.id
	cidr_block              = element(var.private_subnet_cidrs, count.index)
	availability_zone       = element(var.private_availability_zones, count.index)
	count                   = length(var.private_subnet_cidrs)
	tags                    = var.common_tags
}

resource "aws_nat_gateway" "main" {
	count         = length(var.public_subnet_cidrs)
  allocation_id = element(aws_eip.nat.*.id, count.index)
  subnet_id     = element(aws_subnet.public.*.id, count.index)
  depends_on    = [aws_internet_gateway.main]
	tags          = var.common_tags
}

resource "aws_eip" "nat" {
	count = length(var.public_subnet_cidrs)
	vpc   = true
	tags  = var.common_tags
}

resource "aws_route_table" "private" {
	count  = length(var.private_subnet_cidrs)
	vpc_id = aws_vpc.main.id
	tags   = var.common_tags
}

resource "aws_route" "private" {
	count                  = length(compact(var.private_subnet_cidrs))
  route_table_id         = element(aws_route_table.private.*.id, count.index)
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = element(aws_nat_gateway.main.*.id, count.index)
}

resource "aws_route_table_association" "private" {
  count          = length(var.private_subnet_cidrs)
  subnet_id      = element(aws_subnet.private.*.id, count.index)
  route_table_id = element(aws_route_table.private.*.id, count.index)
}
