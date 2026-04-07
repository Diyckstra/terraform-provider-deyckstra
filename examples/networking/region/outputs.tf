output "vpc_id" {
  value = aws_vpc.main.id
}

output "primary_subnet_id" {
  value = module.primary_subnet.subnet_id
}

output "secondary_subnet_id" {
  value = length(module.secondary_subnet) > 0 ? module.secondary_subnet[0].subnet_id : "no secondary subnet"
}
