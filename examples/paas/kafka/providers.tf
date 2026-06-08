terraform {
  required_version = ">= 0.13"

  required_providers {
    aws = {
      source  = "c2devel/rockitcloud"
      version = "~> 25.4"
    }
  }
}

provider "aws" {
  # For K2 Cloud, specify one of the supported regions
  # (e.g. "ru-msk"). For other cloud platforms, enter a
  # non-empty string (e.g. "region-1") and override the
  # API endpoints via environment variables.
  region = var.region
}
