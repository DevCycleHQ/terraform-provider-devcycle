package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"devcycle": func() (tfprotov6.ProviderServer, error) {
		return tfsdk.NewProtocol6Server(New("0.0.3")()), nil
	},
}

func testAccPreCheck(t *testing.T) {
	t.Setenv("DEVCYCLE_CLIENT_ID", os.Getenv("DEVCYCLE_CLIENT_ID"))
	t.Setenv("DEVCYCLE_CLIENT_SECRET", os.Getenv("DEVCYCLE_CLIENT_SECRET"))
	t.Setenv("DEVCYCLE_ACCESS_TOKEN", os.Getenv("DEVCYCLE_ACCESS_TOKEN"))
	t.Setenv("DEVCYCLE_SERVER_TOKEN", os.Getenv("DEVCYCLE_SERVER_TOKEN"))
}
