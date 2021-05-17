resource "random_string" "postgres_password" {
  length = 30
  special = false
}

resource "aws_db_subnet_group" "postgres_subnet_group" {
  name       = "${var.service_name}-${var.environment}-db-subnet-group"
  subnet_ids = split(",", data.aws_ssm_parameter.private_subnets_ssm.value)
}

resource "aws_db_instance" "postgres" {
  allocated_storage      = var.rds_allocated_storage
  engine                 = "postgres"
  engine_version         = "12.6"
  instance_class         = var.rds_instance_class
  identifier             = "${var.service_name}-${var.environment}"
  name                   = var.rds_database_name
  username               = "postgres"
  password               = random_string.postgres_password.result
  parameter_group_name   = "default.postgres12"
  skip_final_snapshot    = true
  db_subnet_group_name   = aws_db_subnet_group.postgres_subnet_group.name
  vpc_security_group_ids = [data.aws_ssm_parameter.vpc_default_security_group_ssm.value]

  depends_on = [
    aws_db_subnet_group.postgres_subnet_group,
  ]
}

resource "aws_ssm_parameter" "postgres_db_conn_url_ssm" {
  name  = "/${var.service_name}/${var.environment}/DB_CONN_URL"
  type  = "SecureString"
  value = "postgres://${aws_db_instance.postgres.username}:${random_string.postgres_password.result}@${aws_db_instance.postgres.endpoint}/${var.rds_database_name}"

  depends_on = [
    aws_db_instance.postgres,
  ]
}
