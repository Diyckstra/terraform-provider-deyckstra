variable "region" {
  description = "The region to create the infrastructure in."
  default     = "ru-msk"
}

variable "first_access_key" {
  type        = string
  description = "First account access key."
}

variable "first_secret_key" {
  type        = string
  description = "First account secret key."
}

variable "second_access_key" {
  type        = string
  description = "Second account access key."
}

variable "second_secret_key" {
  type        = string
  description = "Second account secret key."
}

variable "second_account_id" {
  type        = string
  description = "Second account id (project@customer)."
}
