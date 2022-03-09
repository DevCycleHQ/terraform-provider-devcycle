package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type featureDataSourceType struct{}

func (t featureDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DevCycle Feature data source.",

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
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"key": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Variation key",
					},
					"name": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Variation name",
					},
					"variables": {
						Type:                types.MapType{ElemType: types.StringType},
						Required:            true,
						MarkdownDescription: "Variation variables",
					},
					"id": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variation type",
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"variables": {
				MarkdownDescription: "Feature variables",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:                types.StringType,
						Optional:            true,
						MarkdownDescription: "Variation name",
					},
					"description": {
						Type:                types.StringType,
						Optional:            true,
						MarkdownDescription: "Variation feature key",
					},
					"key": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Variation key",
					},
					"feature_key": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variation feature key",
					},
					"type": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Variation type",
					},
					"id": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Variation type",
					},
					"created_at": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Created at timestamp",
					},
					"updated_at": {
						Type:                types.StringType,
						Computed:            true,
						MarkdownDescription: "Updated at timestamp",
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

	feature, httpResponse, err := d.provider.MgmtClient.FeaturesApi.FeaturesControllerFindOne(ctx, data.Key.Value, data.ProjectKey.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Id = types.String{Value: feature.Id}
	data.Name = types.String{Value: feature.Name}
	data.Key = types.String{Value: feature.Key}
	data.Description = types.String{Value: feature.Description}
	data.ProjectId = types.String{Value: feature.Project}
	data.ProjectKey = types.String{Value: feature.Project}
	data.Type = types.String{Value: feature.Type_}
	data.Variables = variableToTF(feature.Variables)
	data.Variations = variationToTF(feature.Variations, data.Variables)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
