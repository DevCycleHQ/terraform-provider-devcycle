package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	testAccPreCheck(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 nil,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProjectResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("devcycle_project.test", "key", testAccProjectResourceKey()),
					resource.TestCheckResourceAttr("devcycle_project.test", "description", "Terraform acceptance testing"),
				),
			},
			{
				Config: testAccProjectResourceConfigEdit(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("devcycle_project.test", "description", "Terraform acceptance testing-edit"),
				),
			},
			{
				Config:  testAccProjectResourceConfig(),
				Destroy: true,
			},
		},
	})
}

func testAccProjectResourceKey() string {
	return "terraform-acceptance-testing" + randString
}

func testAccProjectResourceConfig() string {
	return `
resource "devcycle_project" "test" {
  name = "TerraformAccTest` + randString + `"
  key = "` + testAccProjectResourceKey() + `"
  description = "Terraform acceptance testing"
}
`
}

func testAccProjectResourceConfigEdit() string {
	return `
resource "devcycle_project" "test" {
  name = "TerraformAccTest` + randString + `"
  key = "` + testAccProjectResourceKey() + `"
  description = "Terraform acceptance testing-edit"
}
`
}
