package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVariableDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccVariableDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devcycle_variable.test", "id", "622117604b06357d06d1d0f9"),
				),
			},
		},
	})
}

const testAccVariableDataSourceConfig = `
data "devcycle_variable" "test" {
  key = "terraform-provider-variable"
  project_key = "terraform-provider-testing"
}
`
