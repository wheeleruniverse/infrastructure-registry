
provider "google" {
  project = var.project
  region  = "us-central1"
}

terraform {
  required_version = ">= 0.14.9"

  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "wheelers-websites"

    workspaces {
      name = "resume-gcp"
    }
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 3.80.0"
    }
  }
}
