resource "random_string" "jwt_secret" {
  length           = 16
  special          = true
  override_special = "/@£$"
}

resource "aws_ssm_parameter" "jwt_secret_ssm" {
  name  = "/${var.service_name}/${var.environment}/JWT_SECRET"
  type  = "SecureString"
  value = random_string.jwt_secret.result
}
