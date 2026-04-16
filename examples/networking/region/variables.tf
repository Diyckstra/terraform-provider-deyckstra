variable "region" {
  description = "The region to set up a network within."
  type        = string
}

variable "base_cidr_block" {
  type = string
}

provider "aws" {
  # For K2 Cloud, specify one of the supported regions.
  # For other cloud platforms, enter a non-empty string,
  # for example, "region-1", and API endpoints.
  region = var.region
}
