terraform {
  required_providers {
    database_dsn = {
      source = "elasticdotventures/database_dsn"
    }
  }
}

provider "database_dsn" {}

data "database_dsn_build" "postgres" {
  driver   = "postgres"
  user     = "myuser"
  password = var.db_password
  host     = "localhost"
  port     = 5432
  name     = "myapp"
  params = {
    sslmode = "require"
  }
}

data "database_dsn_build" "mysql" {
  driver   = "mysql"
  user     = "root"
  password = var.mysql_password
  host     = "mysql.example.com"
  port     = 3306
  name     = "production"
  params = {
    charset = "utf8mb4"
    parseTime = "True"
    loc = "Local"
  }
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "mysql_password" {
  description = "MySQL password"
  type        = string
  sensitive   = true
}

output "postgres_dsn" {
  value     = data.database_dsn_build.postgres.dsn
  sensitive = true
}

output "mysql_dsn" {
  value     = data.database_dsn_build.mysql.dsn
  sensitive = true
}