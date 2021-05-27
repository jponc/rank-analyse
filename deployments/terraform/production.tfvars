environment  = "production"
service_name = "rank-analyse"
aws_profile  = "uptactics"
aws_region   = "us-east-1"

# RDS
rds_instance_class    = "db.t3.small"
rds_allocated_storage = 40
rds_database_name     = "mydatabase"
rds_max_connections   = 500

# VPC
private_subnets_ssm_path            = "/uptactics/production/PRIVATE_SUBNET_IDS"
vpc_default_security_group_ssm_path = "/uptactics/production/DEFAULT_SECURITY_GROUP"
