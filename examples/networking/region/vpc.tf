resource "aws_vpc" "main" {
  cidr_block = cidrsubnet(var.base_cidr_block, 4, 1)

  tags = {
    Name = "terraform-networking-example"
  }
}
