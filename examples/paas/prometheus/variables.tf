variable "ssh_public_key" {
  description = "SSH public key content for the key pair"
  type        = string
}

variable "region" {
  description = "K2 Cloud region"
  type        = string
  default     = "ru-msk"
}

variable "availability_zone" {
  description = "Availability zone for the subnet (match your region, e.g. ru-msk-vol52)"
  type        = string
  default     = "ru-msk-vol52"
}

variable "paas_service_name" {
  description = "PaaS service name (must be unique in the account)"
  type        = string
  default     = "tf-prometheus"
}

variable "instance_type" {
  description = "Instance type for Prometheus nodes"
  type        = string
  default     = "c5.large"
}

variable "webhook_url" {
  description = "Webhook URL for the notification channel; provide an endpoint controlled by you or an intentionally non-routable test URL"
  type        = string
}
