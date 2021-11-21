
provider "aws" {
  region = var.region
}

terraform {
  required_version = ">= 0.14.9"

  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "wheelers-websites"

    workspaces {
      name = "resume-aws"
    }
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 3.55.0"
    }
  }
}
