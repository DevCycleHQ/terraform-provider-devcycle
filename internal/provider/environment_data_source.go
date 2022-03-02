package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type environmentDataSourceType struct{}

func (t environmentDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]tfsdk.Attribute{
			"key": {
				MarkdownDescription: "Project key, usually the lowercase, kebab case name of the project",
				Required:            true,
				Type:                types.StringType,
			},
			"project_id": {
				MarkdownDescription: "Project ID",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Project Id",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t environmentDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return environmentDataSource{
		provider: provider,
	}, diags
}

type environmentDataSourceData struct {
	Id          types.String                    `tfsdk:"id"`
	Key         types.String                    `tfsdk:"key"`
	Name        types.String                    `tfsdk:"name"`
	Description types.String                    `tfsdk:"description"`
	Color       types.String                    `tfsdk:"color"`
	Type        types.String                    `tfsdk:"type"`
	Settings    environmentResourceDataSettings `tfsdk:"settings"`
	ProjectId   types.String                    `tfsdk:"project_id"`
	SDKKeys     []environmentResourceDataSDKKey `tfsdk:"sdkKeys"`
}

type environmentDataSource struct {
	provider provider
}

func (d environmentDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data environmentDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	environment, httpResponse, err := d.provider.MgmtClient.EnvironmentsApi.EnvironmentsControllerFindOne(ctx, data.Key.Value, data.ProjectId.Value)
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read environment, got error: %s", err))
		return
	}

	data.Id.Value = environment.Id
	data.Key.Value = environment.Key
	data.Name.Value = environment.Name
	data.Description.Value = environment.Description
	data.Color.Value = environment.Color
	data.Type.Value = environment.Type_
	data.ProjectId.Value = environment.Project
	data.Settings.AppIconURI.Value = environment.Settings.AppIconURI
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert("mobile", environment.SdkKeys.Mobile)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert("server", environment.SdkKeys.Server)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert("client", environment.SdkKeys.Client)...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
