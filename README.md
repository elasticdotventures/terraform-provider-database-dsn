# Terraform Provider Database DSN

A Terraform provider for building and parsing database DSN (Data Source Name) connection strings. This provider eliminates HCL toil when handling database connection components like `DB_HOST`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, etc.

## Features

- **`database_dsn_build`**: Build database DSNs from individual components
- **`database_dsn_parse`**: Parse existing DSNs into individual components
- Support for multiple database drivers (PostgreSQL, MySQL, SQL Server, etc.)
- Secure handling of sensitive data (passwords, DSNs)
- Parameter support for additional connection options

## Data Sources

### `database_dsn_build`

Constructs a database DSN from component parts.

```hcl
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

output "connection_string" {
  value     = data.database_dsn_build.postgres.dsn
  sensitive = true
}
```

### `database_dsn_parse`

Parses an existing DSN into component parts.

```hcl
data "database_dsn_parse" "existing" {
  dsn = var.database_url
}

output "db_host" {
  value = data.database_dsn_parse.existing.host
}

output "db_port" {
  value = data.database_dsn_parse.existing.port
}
```

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Installation

This provider is available on the [Terraform Registry](https://registry.terraform.io/providers/elasticdotventures/database_dsn).

```hcl
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
```

## Use Cases

- **Environment Configuration**: Reduce the number of environment variables by storing database components in a single DSN
- **Database Migration**: Parse existing connection strings and reconstruct them with different parameters
- **Multi-Environment Deployments**: Build different DSNs for different environments while keeping the same Terraform configuration
- **Security**: Centralize sensitive database credentials while maintaining component access

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
