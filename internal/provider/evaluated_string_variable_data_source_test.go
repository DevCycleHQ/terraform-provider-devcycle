package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEvaluatedStringFeatureDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEvaluatedStringVariableDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devcycle_evaluated_variable_string.test", "value", "String"),
				),
			},
		},
	})
}

const testAccEvaluatedStringVariableDataSourceConfig = `
data "devcycle_evaluated_variable_string" "test" {
  id = "acceptance-testing-string"
  user = {
	id = "acceptancetesting"
  }
  default_value = false
}
`
