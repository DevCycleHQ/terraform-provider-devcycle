package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProjectResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("devcycle_project.test", "key", testAccProjectResourceKey),
					resource.TestCheckResourceAttr("devcycle_project.test", "description", "Terraform acceptance testing"),
				),
			},
			{
				Config: testAccProjectResourceConfigEdit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("devcycle_project.test", "description", "Terraform acceptance testing-edit"),
				),
			},
			{
				Config:  testAccProjectResourceConfig,
				Destroy: true,
			},
		},
	})
}

var testAccProjectResourceKey = "terraform-acceptance-testing" + randString

var testAccProjectResourceConfig = `
resource "devcycle_project" "test" {
  name = "TerraformAccTest` + randString + `"
  key = "` + testAccProjectResourceKey + `"
  description = "Terraform acceptance testing"
}
`

var testAccProjectResourceConfigEdit = `
resource "devcycle_project" "test" {
  name = "TerraformAccTest` + randString + `"
  key = "` + testAccProjectResourceKey + `"
  description = "Terraform acceptance testing-edit"
}
`
