package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type environmentDataSourceType struct{}

func (t environmentDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `DevCycle Environment Data Source. Read data from a given DevCycle Environment. R`,

		Attributes: map[string]tfsdk.Attribute{
			"project_id": {
				MarkdownDescription: `Project id of the project to which the environment belongs.`,
				Computed:            true,
				Type:                types.StringType,
			},
			"project_key": {
				MarkdownDescription: "Project key or id of the project to which the environment belongs",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Environment Name",
				Computed:            true,
				Type:                types.StringType,
			},
			"key": {
				MarkdownDescription: "Environment Key (Human readable id)",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "Environment Description",
				Computed:            true,
				Type:                types.StringType,
			},
			"color": {
				MarkdownDescription: "Environment Color in Hex with leading #",
				Computed:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Environment Type",
				Computed:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Environment Id",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
			"sdk_keys": {
				Computed:            true,
				MarkdownDescription: "SDK Keys for the environment",
				Type:                types.ListType{ElemType: types.StringType},
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
	Id          types.String `tfsdk:"id"`
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	Type        types.String `tfsdk:"type"`
	ProjectId   types.String `tfsdk:"project_id"`
	ProjectKey  types.String `tfsdk:"project_key"`
	SDKKeys     []string     `tfsdk:"sdk_keys"`
}

type environmentDataSource struct {
	provider provider
}

func (d environmentDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data environmentDataSourceData
	if !d.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. Authentication is required to be configured.",
		)
		return
	}
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	environment, httpResponse, err := d.provider.MgmtClient.EnvironmentsApi.EnvironmentsControllerFindOne(ctx, data.Key.Value, data.ProjectKey.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Id = types.String{Value: environment.Id}
	data.Key = types.String{Value: environment.Key}
	data.Name = types.String{Value: environment.Name}
	data.Description = types.String{Value: environment.Description}
	data.Color = types.String{Value: environment.Color}
	data.Type = types.String{Value: environment.Type_}
	data.ProjectId = types.String{Value: environment.Project}
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Mobile)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Server)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Client)...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
