data "aws_availability_zones" "all" {}

module "primary_subnet" {
  source            = "../subnet"
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.all.names[0]
  subnet_index      = 0
}

module "secondary_subnet" {
  # Create a secondary subnet if there are enough AZs in the region
  count = length(data.aws_availability_zones.all.names) > 1 ? 1 : 0

  source            = "../subnet"
  vpc_id            = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.all.names[1]
  subnet_index      = 1
}
