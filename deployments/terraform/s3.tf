resource "aws_s3_bucket" "results_bucket" {
  bucket = local.s3_results_bucket_name
  acl    = "private"


  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["*"]
    max_age_seconds = 3000
  }

  tags = {
    Name = local.s3_results_bucket_name
  }
}

resource "aws_ssm_parameter" "s3_results_bucket_name_ssm" {
  name  = "/${var.service_name}/${var.environment}/S3_RESULTS_BUCKET_NAME"
  type  = "String"
  value = local.s3_results_bucket_name

  depends_on = [
    aws_s3_bucket.results_bucket,
  ]
}
