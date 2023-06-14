package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEvaluatedNumberFeatureDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEvaluatedNumberVariableDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devcycle_evaluated_variable_number.test", "value", "69"),
				),
			},
		},
	})
}

const testAccEvaluatedNumberVariableDataSourceConfig = `
data "devcycle_evaluated_variable_number" "test" {
  key = "acceptance-testing-number"
  user = {
	id = "acceptancetesting"
  }
  default_value = 1
}
`
