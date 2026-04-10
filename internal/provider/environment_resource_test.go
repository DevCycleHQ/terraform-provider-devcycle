package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEnvironmentResource(t *testing.T) {
	testAccPreCheck(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 nil,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEnvironmentResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("devcycle_environment.test", "project_id", "622112634cabe0e9fbaf974d"),
				),
			},
			{
				Config:  testAccEnvironmentResourceConfig(),
				Destroy: true,
			},
		},
	})
}

func testAccEnvironmentResourceConfig() string {
	return `
resource "devcycle_environment" "test" {
  project_id = "622112634cabe0e9fbaf974d"
  name = "TerraformAccTest` + randString + `"
  key = "terraform-acceptance-testing` + randString + `"
  description = "Terraform acceptance testing"
  color = "#232323"
  type = "development"
  settings = {
	app_icon_uri = "test"
  }
}
`
}
