variable "region" {
  description = "The region to set up a network within."
  type        = string
}

variable "base_cidr_block" {
  type = string
}

variable "access_key" {
  description = "The access key taken from c2rc.sh."
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key taken from c2rc.sh."
  type        = string
  sensitive   = true
}
