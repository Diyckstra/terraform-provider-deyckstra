resource "aws_subnet" "main" {
  cidr_block        = cidrsubnet(data.aws_vpc.target.cidr_block, 2, var.subnet_index)
  vpc_id            = var.vpc_id
  availability_zone = var.availability_zone

  tags = {
    Name = "terraform-networking-example"
  }
}

resource "aws_route_table" "main" {
  vpc_id = var.vpc_id

  tags = {
    Name = "terraform-networking-example"
  }
}

resource "aws_route_table_association" "main" {
  subnet_id      = aws_subnet.main.id
  route_table_id = aws_route_table.main.id
}
