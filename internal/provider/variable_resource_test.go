package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVariableResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccVariableResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("devcycle_variable.test", "key", testAccVariableResourceKey),
				),
			},
			{
				Config:  testAccVariableResourceConfig,
				Destroy: true,
			},
		},
	})
}

var testAccVariableResourceKey = "terraform-acceptance-testing" + randSeq(5)

var testAccVariableResourceConfig = `
resource "devcycle_variable" "test" {
  name = "TerraformAccTest` + randSeq(5) + `"
  key = "` + testAccVariableResourceKey + `"
  description = "Terraform acceptance testing"
  type = "Boolean"
  feature_id = "622115014b06357d06d1cf3e"
  project_id = "622112634cabe0e9fbaf974d"
  default_value = "false"
}
`
