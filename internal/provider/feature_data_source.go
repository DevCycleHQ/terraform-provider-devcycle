package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type featureDataSourceType struct{}

func (t featureDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Feature name",
				Computed:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "Feature description",
				Computed:            true,
				Type:                types.StringType,
			},
			"key": {
				MarkdownDescription: "Feature key",
				Required:            true,
				Type:                types.StringType,
			},
			"project_key": {
				MarkdownDescription: "Project key that the feature belongs to",
				Required:            true,
				Type:                types.StringType,
			},
			"project_id": {
				MarkdownDescription: "Project ID that the feature belongs to",
				Computed:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Feature Type",
				Computed:            true,
				Type:                types.StringType,
			},
			"variations": {
				MarkdownDescription: "Feature variations",
				Computed:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"key": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variation key",
					},
					"name": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variation name",
					},
					"id": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variation ID",
						PlanModifiers: tfsdk.AttributePlanModifiers{
							tfsdk.RequiresReplace(),
						},
					},
					"variables": {
						Type:                types.MapType{ElemType: types.StringType},
						Computed:            true,
						MarkdownDescription: "Variation variables - force casted to a string because of nested attributes",
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"variables": {
				MarkdownDescription: "Feature variables",
				Computed:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"key": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variable key",
					},
					"name": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variable name",
					},
					"feature_key": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Feature that this variable is attached to",
					},
					"type": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variable datatype",
					},
					"default_string_value": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variable default value if the type is string",
					},
					"default_json_value": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variable default value if the type is json",
					},
					"default_bool_value": {
						Type:                types.BoolType,
						Computed:            true,
						MarkdownDescription: "Variable default value if the type is bool",
					},
					"default_number_value": {
						Type:                types.NumberType,
						Computed:            true,
						MarkdownDescription: "Variable default value if the type is number",
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Feature ID",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t featureDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return featureDataSource{
		provider: provider,
	}, diags
}

type featureDataSourceData struct {
	Id          types.String                   `tfsdk:"id"`
	Name        types.String                   `tfsdk:"name"`
	Key         types.String                   `tfsdk:"key"`
	Description types.String                   `tfsdk:"description"`
	ProjectId   types.String                   `tfsdk:"project_id"`
	ProjectKey  types.String                   `tfsdk:"project_key"`
	Type        types.String                   `tfsdk:"type"`
	Variations  []featureResourceDataVariation `tfsdk:"variations"`
	Variables   []featureResourceDataVariable  `tfsdk:"variables"`
}

type featureDataSource struct {
	provider provider
}

func (d featureDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data featureDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	feature, httpResponse, err := d.provider.MgmtClient.FeaturesApi.FeaturesControllerFindOne(ctx, data.Key.Value, data.ProjectKey.Value)
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read feature, got error: %s", err))
		return
	}

	data.Id = types.String{Value: feature.Id}
	data.Name = types.String{Value: feature.Name}
	data.Key = types.String{Value: feature.Key}
	data.Description = types.String{Value: feature.Description}
	data.ProjectId = types.String{Value: feature.Project}
	data.ProjectKey = types.String{Value: feature.Project}
	data.Type = types.String{Value: feature.Type_}
	data.Variations = variationToTF(feature.Variations)
	data.Variables = variableToTF(feature.Variables)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
