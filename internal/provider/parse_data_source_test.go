// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestParseDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccParseDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.database_dsn_parse.test", "driver", "postgres"),
					resource.TestCheckResourceAttr("data.database_dsn_parse.test", "user", "testuser"),
					resource.TestCheckResourceAttr("data.database_dsn_parse.test", "host", "localhost"),
					resource.TestCheckResourceAttr("data.database_dsn_parse.test", "port", "5432"),
					resource.TestCheckResourceAttr("data.database_dsn_parse.test", "name", "testdb"),
					resource.TestCheckResourceAttr("data.database_dsn_parse.test", "params.sslmode", "disable"),
					resource.TestCheckResourceAttrSet("data.database_dsn_parse.test", "id"),
				),
			},
		},
	})
}

const testAccParseDataSourceConfig = `
data "database_dsn_parse" "test" {
  dsn = "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
}
`