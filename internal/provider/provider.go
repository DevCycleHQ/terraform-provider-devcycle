package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/dvc_oauth"
	"os"

	dvc_mgmt "github.com/devcyclehq/go-mgmt-sdk"
	dvc_server "github.com/devcyclehq/go-server-sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	MgmtClient          *dvc_mgmt.DVCClient
	ServerClient        *dvc_server.DVCClient
	AccessToken         string
	ServerClientContext context.Context

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
	ServerSDKToken types.String `tfsdk:"server_sdk_token"`
	ClientId       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ClientId.Value != "" && data.ClientSecret.Value != "" {
		auth, err := dvc_oauth.GetAuthToken(data.ClientId.Value, data.ClientSecret.Value)
		if err != nil {
			p.configured = false
			return
		}
		p.AccessToken = auth.AccessToken
	} else {
		clientId := os.Getenv("DEVCYCLE_CLIENT_ID")
		clientSecret := os.Getenv("DEVCYCLE_CLIENT_SECRET")
		if clientId != "" && clientSecret != "" {
			auth, err := dvc_oauth.GetAuthToken(clientId, clientSecret)
			if err != nil {
				p.configured = false
				return
			}
			p.AccessToken = auth.AccessToken
		}
	}

	if data.ServerSDKToken.Value != "" {
		p.ServerClientContext = context.WithValue(context.Background(), dvc_server.ContextAPIKey, dvc_server.APIKey{
			Key: data.ServerSDKToken.Value,
		})
	} else {
		p.ServerClientContext = context.WithValue(context.Background(), dvc_server.ContextAPIKey, dvc_server.APIKey{
			Key: os.Getenv("DEVCYCLE_SERVER_TOKEN"),
		})
	}

	config := dvc_mgmt.NewConfiguration()
	config.AddDefaultHeader("Authorization", p.AccessToken)
	config.BasePath = "https://api.devcycle.com"
	config.UserAgent = "terraform-provider-devcycle"
	p.MgmtClient = dvc_mgmt.NewAPIClient(config)
	p.ServerClient = dvc_server.NewDVCClient()
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
		"devcycle_project":                    projectDataSourceType{},
		"devcycle_environment":                environmentDataSourceType{},
		"devcycle_feature":                    featureDataSourceType{},
		"devcycle_variable":                   variableDataSourceType{},
		"devcycle_evaluated_variable_boolean": evaluatedBoolVariableDataSourceType{},
		"devcycle_evaluated_variable_string":  evaluatedStringVariableDataSourceType{},
		"devcycle_evaluated_variable_number":  evaluatedNumberVariableDataSourceType{},
		"devcycle_evaluated_variable_json":    evaluatedJSONVariableDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "This provider allows you to manage DevCycle projects, environments, features, and variables. It uses the DevCycle API to manage these resources.  You can find more information about the DevCycle API [here](https://docs.devcycle.com/management-api/)." +
			"\n\n" + "This provider is compatible with Terraform v1.0 and newer. Because of the way that authentication for the management api works - this provider will have access to manage all projects within a DevCycle org. Be careful!",
		Attributes: map[string]tfsdk.Attribute{
			"client_id": {
				MarkdownDescription: "API Authentication Client ID. Found in your DevCycle account settings.",
				Optional:            true,
				Sensitive:           true,
				Type:                types.StringType,
			},
			"client_secret": {
				MarkdownDescription: "API Authentication Client Secret. Found in your DevCycle account settings.",
				Optional:            true,
				Sensitive:           true,
				Type:                types.StringType,
			},
			"server_sdk_token": {
				Type:                types.StringType,
				MarkdownDescription: "Server SDK Token. This is specific to a given project, and an environment. Used to identify and authenticate server sdk requests to evaluate feature flags.",
				Sensitive:           true,
				Optional:            true,
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
