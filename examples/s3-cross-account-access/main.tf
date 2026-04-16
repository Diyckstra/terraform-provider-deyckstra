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

data "aws_canonical_user_id" "first" {
  provider = aws.first
}

data "aws_canonical_user_id" "second" {
  provider = aws.second
}

resource "aws_s3_bucket" "example" {
  provider = aws.first

  bucket = "terraform-s3-example"

  tags = {
    Name = "terraform-s3-example"
  }
}

resource "aws_s3_bucket_acl" "example" {
  provider = aws.first

  bucket = aws_s3_bucket.example.id

  access_control_policy {
    grant {
      grantee {
        id   = data.aws_canonical_user_id.first.id
        type = "CanonicalUser"
      }
      permission = "FULL_CONTROL"
    }

    grant {
      grantee {
        id   = data.aws_canonical_user_id.second.id
        type = "CanonicalUser"
      }
      permission = "FULL_CONTROL"
    }

    owner {
      id = data.aws_canonical_user_id.first.id
    }
  }
}

resource "aws_s3_object" "first" {
  provider = aws.first

  depends_on = [aws_s3_bucket_acl.example]

  bucket = aws_s3_bucket.example.id
  key    = "object-uploaded-via-first-creds"
  source = "${path.module}/first.txt"

  tags = {
    Name = "terraform-s3-example"
  }
}

resource "aws_s3_object" "second" {
  provider = aws.second

  depends_on = [aws_s3_bucket_acl.example]

  bucket = aws_s3_bucket.example.id
  key    = "object-uploaded-via-second-creds"
  source = "${path.module}/second.txt"

  tags = {
    Name = "terraform-s3-example"
  }
}
