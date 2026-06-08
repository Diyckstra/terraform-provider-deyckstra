variable "region" {
  description = "Region to create the infrastructure in."
  type        = string
  default     = "ru-msk"
}

variable "availability_zone" {
  description = "Availability zone for the subnet (must be inside region)."
  type        = string
  default     = "ru-msk-vol52"
}

variable "service_name" {
  description = "PaaS service name (must be unique in the project)."
  type        = string
  default     = "terraform-kafka-example"
}

variable "kafka_version" {
  description = "Kafka version to deploy."
  type        = string
  default     = "3.7.0"

  validation {
    condition     = contains(["3.6.1", "3.7.0"], var.kafka_version)
    error_message = "kafka_version must be one of: 3.6.1, 3.7.0."
  }
}

variable "instance_type" {
  description = "Instance type for Kafka broker and coordinator nodes."
  type        = string
  default     = "c5gl20.2large"
}

variable "ssh_public_key" {
  description = "SSH public key content for the key pair used to access Kafka nodes."
  type        = string
}
