
provider "google" {
  project = var.project
  region  = lower(var.location_region)
}

terraform {
  required_version = ">= 0.14.9"

  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "wheelers-websites"

    workspaces {
      name = "cloudguruchallenge-2108-core"
    }
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 3.80.0"
    }
  }
}
