
locals {
  tags = {
    environment = var.environment
    owner       = var.owner
    project     = var.project
  }
}

variable "bucket_replication" {
  default     = "LRS"
  description = "bucket replication"
  type        = string
}

variable "bucket_tier" {
  default     = "Standard"
  description = "bucket tier"
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

variable "tables" {
  default     = ["certification", "education", "experience", "project", "skill", "visitor"]
  description = "database tables to create"
  type        = list(string)
}

variable "tls_version" {
  default     = "TLS1_2"
  description = "tls version"
  type        = string
}
