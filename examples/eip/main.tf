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

resource "aws_eip" "example" {
  instance = aws_instance.example.id

  tags = {
    Name = "terraform-eip-example"
  }
}

# Default security group to access
# the instances over SSH and HTTP
resource "aws_security_group" "example" {
  name        = "terraform-eip-example"
  description = "Used in the terraform"

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "terraform-eip-example"
  }
}

data "aws_ami" "selected" {
  most_recent = true

  name_regex = "Ubuntu 24.04 [Cloud Image]*"

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["k2"]
}

resource "aws_instance" "example" {
  instance_type = "c5.large"
  ami           = data.aws_ami.selected.id

  # The name of the SSH keypair created in the cloud console
  key_name = var.key_name

  # Once the instance is created, a remote provisioner is launched on it.
  # In this case, it installs nginx and starts it. By default,
  # this should be on port 80
  user_data = file("userdata.sh")

  vpc_security_group_ids = [aws_security_group.example.id]

  tags = {
    Name = "terraform-eip-example"
  }
}
