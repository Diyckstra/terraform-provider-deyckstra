variable "region" {
  description = "The region to create the infrastructure in."
  default     = "ru-msk"
}

variable "key_name" {
  type        = string
  description = "The name of the SSH keypair."
}
