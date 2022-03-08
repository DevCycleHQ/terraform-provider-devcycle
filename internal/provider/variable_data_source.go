package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type variableDataSourceType struct{}

func (t variableDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DevCycle Variable data source",

		Attributes: map[string]tfsdk.Attribute{
			"key": {
				MarkdownDescription: "Variable key",
				Required:            true,
				Type:                types.StringType,
			},
			"feature_id": {
				MarkdownDescription: "Feature ID",
				Computed:            true,
				Type:                types.StringType,
			},
			"project_id": {
				MarkdownDescription: "Project ID",
				Computed:            true,
				Type:                types.StringType,
			},
			"project_key": {
				MarkdownDescription: "Project key",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"id": {
				MarkdownDescription: "Variable Id",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"name": {
				MarkdownDescription: "Variable name",
				Computed:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "Variable description",
				Computed:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Variable type",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t variableDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return variableDataSource{
		provider: provider,
	}, diags
}

type variableDataSourceData struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Key         types.String `tfsdk:"key"`
	FeatureId   types.String `tfsdk:"feature_id"`
	ProjectId   types.String `tfsdk:"project_id"`
	ProjectKey  types.String `tfsdk:"project_key"`
	Type        types.String `tfsdk:"type"`
	Id          types.String `tfsdk:"id"`
}

type variableDataSource struct {
	provider provider
}

func (d variableDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data variableDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	variable, httpResponse, err := d.provider.MgmtClient.VariablesApi.VariablesControllerFindOne(ctx, data.Key.Value, data.ProjectKey.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}
	data.Id = types.String{Value: variable.Id}
	data.Key = types.String{Value: variable.Key}
	data.Name = types.String{Value: variable.Name}
	data.Description = types.String{Value: variable.Description}
	data.Type = types.String{Value: variable.Type_}
	data.FeatureId = types.String{Value: variable.Feature}
	data.ProjectId = types.String{Value: variable.Project}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
