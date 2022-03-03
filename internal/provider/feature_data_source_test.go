package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFeatureDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccFeatureDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devcycle_feature.test", "id", "622115014b06357d06d1cf3e"),
				),
			},
		},
	})
}

const testAccFeatureDataSourceConfig = `
data "devcycle_feature" "test" {
  key = "terraform-provider-feature"
  project_key = "terraform-provider-testing"
}
`
