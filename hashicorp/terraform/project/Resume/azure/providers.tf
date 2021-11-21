
provider "azurerm" {
  features {}
}

terraform {
  required_version = ">= 0.14.9"

  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "wheelers-websites"

    workspaces {
      name = "resume-azure"
    }
  }
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 2.73.0"
    }
  }
}
