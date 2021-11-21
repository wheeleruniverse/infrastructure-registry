
provider "google" {
  billing_project       = var.project
  project               = var.project
  region                = lower(var.location_region)
  user_project_override = true
}

terraform {
  required_version = ">= 0.14.9"

  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "wheelers-websites"

    workspaces {
      name = "cloudguruchallenge-2108-services"
    }
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 3.80.0"
    }
  }
}
