resource "aws_rds_cluster" "main" {
  cluster_identifier      = "${var.name}-${var.environment}-aurora"
  engine                  = "aurora-postgresql"
  engine_mode             = "provisioned"
  availability_zones      = var.availability_zones
  database_name           = var.database_settings.DB_NAME
  master_username         = var.database_settings.DB_USER
  master_password         = "postgres_password" # TODO: use AWS secret manager later
  vpc_security_group_ids  = [ aws_security_group.aurora.id ]
  db_subnet_group_name    = aws_db_subnet_group.main.name
  skip_final_snapshot     = true
  
  serverlessv2_scaling_configuration {
    max_capacity = 1.0
    min_capacity = 0.5
  }

  tags = var.common_tags
}

resource "aws_db_subnet_group" "main" {
  name       = "${var.name}-${var.environment}-subnet-group"
  subnet_ids = flatten(aws_subnet.public.*.id)
  tags       = var.common_tags
}

resource "aws_rds_cluster_instance" "main" {
  cluster_identifier   = aws_rds_cluster.main.id
  instance_class       = "db.serverless"
  engine               = aws_rds_cluster.main.engine
  engine_version       = aws_rds_cluster.main.engine_version
  # publicly_accessible  = true
  db_subnet_group_name = aws_db_subnet_group.main.name
  tags                 = var.common_tags
}


# TODO: run create_tables.sql in /internal-health for one time