resource "aws_rds_cluster" "main" {
  cluster_identifier      = "${var.name}-${var.environment}-aurora"
  engine                  = "aurora-postgresql"
  engine_mode             = "provisioned"
  availability_zones      = var.availability_zones
  database_name           = var.database_settings.DB_NAME
  master_username         = var.database_settings.DB_USER
  master_password         = "postgres_password" # TODO: use AWS secret manager later
  vpc_security_group_ids  = [ aws_security_group.aurora.id ]
  
  serverlessv2_scaling_configuration {
    max_capacity = 1.0
    min_capacity = 0.5
  }

  tags = var.common_tags
}

resource "aws_rds_cluster_instance" "main" {
  cluster_identifier = aws_rds_cluster.main.id
  instance_class     = "db.serverless"
  engine             = aws_rds_cluster.main.engine
  engine_version     = aws_rds_cluster.main.engine_version
  
  tags               = var.common_tags
}

# TODO: figure out how to set / get database host and port
# TODO: figure out how to run create_tables.sql
# TODO: where to set IAM role (iam_database_authentication_enabled, iam_roles)
# TODO: where to set security groups

# Docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/rds_cluster