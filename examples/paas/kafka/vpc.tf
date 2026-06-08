resource "aws_vpc" "example" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "terraform-paas-kafka-example"
  }
}

resource "aws_subnet" "example" {
  vpc_id            = aws_vpc.example.id
  cidr_block        = cidrsubnet(aws_vpc.example.cidr_block, 4, 1)
  availability_zone = var.availability_zone

  tags = {
    Name = "terraform-paas-kafka-example"
  }
}

# PaaS services require an Internet Gateway in the VPC to reach the platform
# control plane. The platform itself does not expose the brokers publicly.
resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "terraform-paas-kafka-example"
  }
}

resource "aws_route" "example_default" {
  route_table_id         = aws_vpc.example.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.example.id
}
