package main

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/dvc_oauth"
	"os"
)

func main() {
	fmt.Println(dvc_oauth.GetAuthToken(os.Getenv("DEVCYCLE_CLIENT_ID"), os.Getenv("DEVCYCLE_CLIENT_SECRET")))
}
