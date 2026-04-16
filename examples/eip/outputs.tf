output "private_ip" {
  value = aws_instance.example.private_ip
}

output "elastic_ip" {
  value = aws_eip.example.public_ip
}
