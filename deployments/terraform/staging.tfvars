environment  = "staging"
service_name = "rank-analyse"
aws_profile  = "uptactics"
aws_region   = "us-east-1"

# RDS
rds_instance_class    = "db.t3.small"
rds_allocated_storage = 20
rds_database_name     = "mydatabase"
rds_max_connections   = 500

# VPC
private_subnets_ssm_path            = "/uptactics/staging/PRIVATE_SUBNET_IDS"
vpc_default_security_group_ssm_path = "/uptactics/staging/DEFAULT_SECURITY_GROUP"

# Domain
domain_ssm_value = "api-staging-rankanalyze.appextend.com"

# Frontend
frontend_url                       = "rankanalyze-staging.appextend.com"
appextend_hosted_zone_id_ssm_path  = "/uptactics/APPEXTEND_HOSTED_ZONE_ID"
appextend_certificate_arn_ssm_path = "/uptactics/APPEXTEND_CERTIFICATE_ARN"
