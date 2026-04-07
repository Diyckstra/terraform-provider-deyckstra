variable "vpc_id" {
  type = string
}

variable "availability_zone" {
  type = string
}

variable "subnet_index" {
  type        = number
  description = "Unique number used to prevent CIDR overlapping."
}

data "aws_availability_zone" "target" {
  name = var.availability_zone
}

data "aws_vpc" "target" {
  id = var.vpc_id
}
