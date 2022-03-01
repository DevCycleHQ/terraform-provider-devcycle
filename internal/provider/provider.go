package provider

import (
	"context"
	"fmt"

	dvc_mgmt "github.com/devcyclehq/go-mgmt-sdk"
	dvc_server "github.com/devcyclehq/go-server-sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//
	MgmtClient   *dvc_mgmt.APIClient
	ServerClient *dvc_server.DVCClient

	MgmtClientToken   string
	ServerClientToken string

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	ServerSDKToken  types.String `tfsdk:"server_sdk_token"`
	ManagementToken types.String `tfsdk:"management_token"`
	ProjectId       types.String `tfsdk:"project_id"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	p.MgmtClient = dvc_mgmt.NewAPIClient(dvc_mgmt.NewConfiguration())
	p.ServerClient = dvc_server.NewDVCClient()
	p.MgmtClientToken = data.ManagementToken.Value
	p.ServerClientToken = data.ServerSDKToken.Value

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"devcycle_project":     projectResourceType{},
		"devcycle_environment": environmentResourceType{},
		"devcycle_feature":     featureResourceType{},
		"devcycle_variable":    variableResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"scaffolding_example": exampleDataSourceType{},
	}, nil
}

// TODO: Pass in provider config values for configuration of the SDKs.
func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example": {
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"server": {
				Type:                nil,
				Attributes:          nil,
				Description:         "",
				MarkdownDescription: "",
				Required:            false,
				Optional:            false,
				Computed:            false,
				Sensitive:           false,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
