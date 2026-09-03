data "aws_availability_zones" "all" {}

resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-eks-example"
  }
}

resource "aws_subnet" "example" {
  count = length(data.aws_availability_zones.all.names)

  availability_zone = data.aws_availability_zones.all.names[count.index]
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = aws_vpc.example.id

  tags = {
    Name = "terraform-eks-example"
  }
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "terraform-eks-example"
  }
}

resource "aws_nat_gateway" "example" {
  depends_on = [aws_internet_gateway.example]

  vpc_id = aws_vpc.example.id

  tags = {
    Name = "terraform-eks-example"
  }
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.example.id
  }

  tags = {
    Name = "terraform-eks-example"
  }
}

resource "aws_route_table_association" "example" {
  count = length(data.aws_availability_zones.all.names)

  subnet_id      = aws_subnet.example[count.index].id
  route_table_id = aws_route_table.example.id
}
