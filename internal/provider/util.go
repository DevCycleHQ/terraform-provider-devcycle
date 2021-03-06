package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/rand"
	"net/http"
)

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func userDataSchema() tfsdk.Attribute {
	return tfsdk.Attribute{
		MarkdownDescription: "User data to drive bucketing into variations for feature flag evaluations.",
		Required:            true,
		Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "User ID",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "User name",
				Optional:            true,
				Type:                types.StringType,
			},
			"app_version": {
				MarkdownDescription: "User app version",
				Optional:            true,
				Type:                types.StringType,
			},
			"email": {
				MarkdownDescription: "User email",
				Optional:            true,
				Type:                types.StringType,
			},
			"app_build": {
				MarkdownDescription: "User app build",
				Optional:            true,
				Type:                types.StringType,
			},
		}),
		PlanModifiers: tfsdk.AttributePlanModifiers{
			tfsdk.RequiresReplace(),
		},
	}
}

func handleDevCycleHTTP(err error, httpResponse *http.Response, resp *diag.Diagnostics) bool {
	if err != nil || (httpResponse.StatusCode > 299 || httpResponse.StatusCode < 200) {
		resp.AddError("Client Error", fmt.Sprintf("DevCycle Terraform Error: %s.\nHTTP Response: %v", err, httpResponse.Request))
		return true
	}
	return false
}

type evaluatedVariableDataSourceDataUser struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	AppVersion types.String `tfsdk:"app_version"`
	Email      types.String `tfsdk:"email"`
	AppBuild   types.String `tfsdk:"app_build"`
}
