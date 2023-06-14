package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEvaluatedBooleanFeatureDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEvaluatedBoolVariableDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devcycle_evaluated_variable_boolean.test", "value", "true"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEvaluatedBoolVariableDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devcycle_evaluated_variable_boolean.test-default", "value", "false"),
				),
			},
		},
	})
}

const testAccEvaluatedBoolVariableDataSourceConfig = `
data "devcycle_evaluated_variable_boolean" "test" {
  key = "acceptance-testing-boolean"
  user = {
	id = "acceptancetesting"
  }
  default_value = false
}
data "devcycle_evaluated_variable_boolean" "test-default" {
  key = "acceptance-testing-boolean-no-exist"
  user = {
	id = "acceptancetesting"
  }
  default_value = false
}
`
