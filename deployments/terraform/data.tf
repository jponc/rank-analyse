data "aws_ssm_parameter" "private_subnets_ssm" {
  name = var.private_subnets_ssm_path
}

data "aws_ssm_parameter" "vpc_default_security_group_ssm" {
  name = var.vpc_default_security_group_ssm_path
}
