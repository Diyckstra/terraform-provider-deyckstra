terraform {
  required_version = ">= 0.12"

  required_providers {
    aws = {
      source  = "c2devel/rockitcloud"
      version = "~> 25.4"
    }
  }
}

provider "aws" {
  # For K2 Cloud, specify one of the supported regions.
  # For other cloud platforms, enter a non-empty string,
  # for example, "region-1", and API endpoints.
  region = var.region
}

data "aws_ami" "selected" {
  most_recent = true

  name_regex = "Ubuntu*"

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["k2"]
}

resource "aws_instance" "example" {
  instance_type = "c5.large"
  ami           = data.aws_ami.selected.id

  tags = {
    Name = "terraform-count-example"
  }

  # This will create 4 instances
  count = 4
}

