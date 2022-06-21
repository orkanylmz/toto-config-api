variable "project" {}
variable "user" {}
variable "region" {}
variable "zone" {}

variable "billing_account" {
  description = "Billing account display name"
}

variable "repository_name" {
  default = "toto-config-api"
}

