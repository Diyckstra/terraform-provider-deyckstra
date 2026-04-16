output "ids" {
  value = "Instances: ${join(", ", aws_instance.example.*.id)}"
}
