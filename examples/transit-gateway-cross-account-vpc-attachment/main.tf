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
  alias = "first"

  # For K2 Cloud, specify one of the supported regions.
  # For other cloud platforms, enter a non-empty string,
  # for example, "region-1", and API endpoints.
  region     = var.region
  access_key = var.first_access_key
  secret_key = var.first_secret_key
}

provider "aws" {
  alias = "second"

  # For K2 Cloud, specify one of the supported regions.
  # For other cloud platforms, enter a non-empty string,
  # for example, "region-1", and API endpoints.
  region     = var.region
  access_key = var.second_access_key
  secret_key = var.second_secret_key
}

resource "aws_ec2_transit_gateway" "example" {
  provider = aws.first

  tags = {
    Name = "terraform-tgw-example"
  }
}

resource "aws_ec2_transit_gateway_project_access" "example" {
  provider = aws.first

  transit_gateway_id = aws_ec2_transit_gateway.example.id
  account_id = var.second_account_id
}

resource "aws_vpc" "example" {
  provider = aws.second

  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-tgw-example"
  }
}

data "aws_availability_zones" "all" {
  provider = aws.second
}

resource "aws_subnet" "example" {
  provider = aws.second

  availability_zone = data.aws_availability_zones.all.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.example.id

  tags = {
    Name = "terraform-tgw-example"
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "example" {
  provider = aws.second

  depends_on = [aws_ec2_transit_gateway_project_access.example]

  subnet_ids         = [aws_subnet.example.id]
  transit_gateway_id = aws_ec2_transit_gateway.example.id
  vpc_id             = aws_vpc.example.id

  tags = {
    Name = "terraform-tgw-example"
  }
}
