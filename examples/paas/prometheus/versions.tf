terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "c2devel/rockitcloud"
      version = "~> 25.4"
    }
  }
}
