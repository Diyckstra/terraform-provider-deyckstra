resource "aws_security_group" "region" {
  name        = "terraform-networking-example-region"
  description = "Open access within this region"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = [aws_vpc.main.cidr_block]
  }

  tags = {
    Name = "terraform-networking-example"
  }
}

resource "aws_security_group" "internal-all" {
  name        = "terraform-networking-example-internal-all"
  description = "Open access within the full internal network"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = [var.base_cidr_block]
  }

  tags = {
    Name = "terraform-networking-example"
  }
}
