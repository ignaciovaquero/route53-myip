variable "region" {
  type        = string
  description = "The AWS region"
  default     = "eu-south-2"
}

variable "record_conditions" {
  type        = map(string)
  description = "The record conditions"
  default = {
    type   = "A"
    name   = "home.ignaciovaquero.es"
    action = "UPSERT"
  }
}
