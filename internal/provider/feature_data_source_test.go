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
					resource.TestCheckResourceAttr("data.devcycle_feature.test", "project_id", "621fc3113bb541e45c20e6da"),
				),
			},
		},
	})
}

const testAccFeatureDataSourceConfig = `
data "devcycle_feature" "test" {
  key = "experiment-feature"
  project_key = "terraform-provider-testing"
}
`
