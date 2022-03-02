package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEnvironmentDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEnvironmentDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devcycle_environment.test", "key", "example-id"),
				),
			},
		},
	})
}

const testAccEnvironmentDataSourceConfig = `
data "devcycle_environment" "test" {
  key = "example-id"
  project_id = "example-project-id"
}
`
