package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEvaluatedJSONFeatureDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEvaluatedJSONVariableDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devcycle_evaluated_variable_json.test", "value", "{\"object\":true}"),
				),
			},
		},
	})
}

const testAccEvaluatedJSONVariableDataSourceConfig = `
data "devcycle_evaluated_variable_json" "test" {
  key = "acceptance-testing-json"
  user = {
	id = "acceptancetesting"
  }
  default_value = "{}"
}
`
