package provider

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"devcycle": func() (tfprotov6.ProviderServer, error) {
		return tfsdk.NewProtocol6Server(New("testing")()), nil
	},
}
var randString = randSeq(5)

func testAccPreCheck(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Setenv("DEVCYCLE_CLIENT_ID", os.Getenv("DEVCYCLE_CLIENT_ID"))
	t.Setenv("DEVCYCLE_CLIENT_SECRET", os.Getenv("DEVCYCLE_CLIENT_SECRET"))
	t.Setenv("DEVCYCLE_ACCESS_TOKEN", os.Getenv("DEVCYCLE_ACCESS_TOKEN"))
	t.Setenv("DEVCYCLE_SERVER_TOKEN", os.Getenv("DEVCYCLE_SERVER_TOKEN"))
}
