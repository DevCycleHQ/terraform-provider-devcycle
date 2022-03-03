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
				),
			},
			{
				Config:  testAccProjectResourceConfig,
				Destroy: true,
			},
		},
	})
}

var testAccProjectResourceKey = "terraform-acceptance-testing" + randSeq(5)

var testAccProjectResourceConfig = `
resource "devcycle_project" "test" {
  name = "TerraformAccTest` + randSeq(5) + `"
  key = "` + testAccProjectResourceKey + `"
  description = "Terraform acceptance testing"
}
`
