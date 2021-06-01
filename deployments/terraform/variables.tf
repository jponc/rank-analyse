variable "environment" {
  default     = "staging"
  description = "Environment"
}

variable "service_name" {
  description = "This is the name of the service"
}

variable "aws_profile" {
  description = "AWS Profile"
}

variable "aws_region" {
  description = "AWS Region"
}

variable "rds_instance_class" {
  description = "RDS Instance type"
}

variable "rds_allocated_storage" {
  type        = number
  description = "Allocated storage for RDS instance"
}

variable "rds_database_name" {
  type        = string
  description = "RDS database name"
}

variable "rds_max_connections" {
  type        = number
  description = "RDS database name"
}

variable "private_subnets_ssm_path" {
  type        = string
  description = "SSM Path to Private subnets StringList"
}

variable "domain_ssm_value" {
  type        = string
  description = "SSM Domain Value"
}

variable "vpc_default_security_group_ssm_path" {
  type        = string
  description = "SSM Path to VPC default security group"
}
