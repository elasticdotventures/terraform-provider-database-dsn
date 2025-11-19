// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestBuildDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBuildDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.database_dsn_build.test", "driver", "postgres"),
					resource.TestCheckResourceAttr("data.database_dsn_build.test", "user", "testuser"),
					resource.TestCheckResourceAttr("data.database_dsn_build.test", "host", "localhost"),
					resource.TestCheckResourceAttr("data.database_dsn_build.test", "port", "5432"),
					resource.TestCheckResourceAttr("data.database_dsn_build.test", "name", "testdb"),
					resource.TestCheckResourceAttr("data.database_dsn_build.test", "dsn", "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"),
					resource.TestCheckResourceAttrSet("data.database_dsn_build.test", "id"),
				),
			},
		},
	})
}

const testAccBuildDataSourceConfig = `
data "database_dsn_build" "test" {
  driver   = "postgres"
  user     = "testuser"
  password = "testpass"
  host     = "localhost"
  port     = 5432
  name     = "testdb"
  params = {
    sslmode = "disable"
  }
}
`