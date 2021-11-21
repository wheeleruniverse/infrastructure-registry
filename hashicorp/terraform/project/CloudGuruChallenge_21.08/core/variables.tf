
locals {
  labels = {
    environment = var.environment
    project     = var.project
  }
}

variable "dns_zone" {
  default     = "wheelersadvice-com"
  description = "cloud dns managed zone name"
  type        = string
}

variable "domain" {
  default     = "wheelersadvice"
  description = "domain name"
  type        = string
}

variable "environment" {
  default     = "prd"
  description = "environment name"
  type        = string
}

variable "fqdn" {
  default     = "wheelersadvice.com"
  description = "fully qualified domain name"
  type        = string
}

variable "image" {
  description = "container image"
  type        = string
}

variable "location_dual_region" {
  default     = "NAM4"
  description = "a specific pair of regions, such as Tokyo and Osaka"
  type        = string

  validation {
    condition     = contains(["NAM4"], var.location_dual_region)
    error_message = "Variable location_dual_region is invalid."
  }
}

variable "location_multi_region" {
  default     = "US"
  description = "a large geographic area, such as the United States, that contains two or more geographic places"
  type        = string

  validation {
    condition     = contains(["US"], var.location_multi_region)
    error_message = "Variable location_multi_region is invalid."
  }
}

variable "location_region" {
  default     = "US-CENTRAL1"
  description = "a specific geographic place, such as São Paulo"
  type        = string

  validation {
    condition     = contains(["US-CENTRAL1", "US-EAST1", "US-EAST4", "US-WEST1", "US-WEST2", "US-WEST3", "US-WEST4"], var.location_region)
    error_message = "Variable location_region is invalid."
  }
}

variable "project" {
  default     = "cloudguruchallenge-2108"
  description = "project name"
  type        = string
}

variable "service_account_json" {
  description = "service account json"
  type        = string
}

variable "storage_class" {
  default     = "STANDARD"
  description = "storage class"
  type        = string
}
