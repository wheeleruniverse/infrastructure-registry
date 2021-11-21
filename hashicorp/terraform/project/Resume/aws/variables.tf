
locals {
  aws_fqdn_list = [
    "aws.${var.fqdn}",
    var.fqdn,
    "www.${var.fqdn}"
  ]
  tags = {
    environment = var.environment
    owner       = var.owner
    project     = var.project
  }
}

variable "account" {
  default     = "778263278211"
  description = "account id"
  type        = string
}

variable "domain" {
  default     = "wheelercloudguru"
  description = "domain name"
  type        = string
}

variable "environment" {
  default     = "prd"
  description = "environment name"
  type        = string
}

variable "fqdn" {
  default     = "wheelercloudguru.com"
  description = "fully qualified domain name"
  type        = string
}

variable "image" {
  description = "container image to deploy"
  type        = string
}

variable "lambda_role" {
  default     = "arn:aws:iam::778263278211:role/wheelercloudguru-lambda"
  description = "iam role to assign to the lambda function"
  type        = string
}

variable "owner" {
  default     = "justin.wheeler@wheelerswebsites.com"
  description = "owner email address"
  type        = string
}

variable "project" {
  default     = "resume"
  description = "project name"
  type        = string
}

variable "region" {
  default     = "us-east-1"
  description = "region name"
  type        = string
}

variable "tables" {
  default     = ["certification", "education", "experience", "project", "skill", "visitor"]
  description = "database tables to create"
  type        = list(string)
}

variable "tls_version" {
  default     = "TLSv1.2_2021"
  description = "tls version"
  type        = string
}
