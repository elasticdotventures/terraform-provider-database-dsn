terraform {
  required_providers {
    database_dsn = {
      source = "elasticdotventures/database_dsn"
    }
  }
}

provider "database_dsn" {
  # No configuration required
}
