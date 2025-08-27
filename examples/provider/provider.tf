terraform {
  required_providers {
    database_dsn = {
      source = "elasticdotventures/database-dsn"
    }
  }
}

provider "database_dsn" {
  # No configuration required
}
