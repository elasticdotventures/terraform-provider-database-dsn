terraform {
  required_providers {
    database_dsn = {
      source = "elasticdotventures/database_dsn"
    }
  }
}

provider "database_dsn" {}

data "database_dsn_parse" "existing" {
  dsn = var.database_url
}

variable "database_url" {
  description = "Existing database URL to parse"
  type        = string
  sensitive   = true
}

output "parsed_driver" {
  value = data.database_dsn_parse.existing.driver
}

output "parsed_host" {
  value = data.database_dsn_parse.existing.host
}

output "parsed_port" {
  value = data.database_dsn_parse.existing.port
}

output "parsed_database" {
  value = data.database_dsn_parse.existing.name
}

output "parsed_user" {
  value = data.database_dsn_parse.existing.user
}

output "parsed_params" {
  value = data.database_dsn_parse.existing.params
}