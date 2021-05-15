terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  backend "s3" {
    bucket         = "uptactics-terraform-up-and-running-state"
    key            = "global/s3/terraform.tfstate"
    region         = "us-east-1"
    profile        = "uptactics"
    dynamodb_table = "uptactics-terraform-up-and-running-locks"
    encrypt        = true
  }
}

provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
}

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
