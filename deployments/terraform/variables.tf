variable "environment" {
  default = "staging"
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
