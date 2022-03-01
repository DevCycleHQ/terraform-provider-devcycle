package provider

import (
	"context"
	"fmt"
	devcyclem "github.com/devcyclehq/go-mgmt-sdk"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type featureResourceType struct{}

func (t featureResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Feature name",
				Required:            true,
				Type:                types.StringType,
			},

			"description": {
				MarkdownDescription: "Feature description",
				Required:            true,
				Type:                types.StringType,
			},

			"key": {
				MarkdownDescription: "Feature key",
				Required:            true,
				Type:                types.StringType,
			},
			"projectId": {
				MarkdownDescription: "Project ID that the feature belongs to",
				Required:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Feature Type",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "Feature tags",
				Required:            true,
				Type:                types.ListType{ElemType: types.StringType},
			},
			"variations": {
				MarkdownDescription: "Feature variations",
				Required:            true,
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
					"description": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Variation description",
					},
					"variables": {
						Type:                types.MapType{ElemType: types.ObjectType{}},
						Required:            true,
						MarkdownDescription: "Variation variables",
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"variables": {
				MarkdownDescription: "Feature variables",
				Required:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"key": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Variable key",
					},
					"name": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Variable name",
					},
					"feature_key": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Feature that this variable is attached to",
					},
					"type": {
						Type:                types.StringType,
						Required:            true,
						MarkdownDescription: "Variable datatype",
					},
					"default_string_value": {
						Type:                types.StringType,
						Optional:            true,
						MarkdownDescription: "Variable default value if the type is string",
					},
					"default_json_value": {
						Type:                types.StringType,
						Optional:            true,
						MarkdownDescription: "Variable default value if the type is json",
					},
					"default_bool_value": {
						Type:                types.BoolType,
						Optional:            true,
						MarkdownDescription: "Variable default value if the type is bool",
					},
					"default_number_value": {
						Type:                types.NumberType,
						Optional:            true,
						MarkdownDescription: "Variable default value if the type is number",
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Feature ID",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t featureResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return featureResource{
		provider: provider,
	}, diags
}

type featureResourceData struct {
	Id          types.String                   `tfsdk:"id"`
	Name        types.String                   `tfsdk:"name"`
	Key         types.String                   `tfsdk:"key"`
	Description types.String                   `tfsdk:"description"`
	ProjectId   types.String                   `tfsdk:"projectId"`
	Source      types.String                   `tfsdk:"source"`
	Type        types.String                   `tfsdk:"type"`
	Tags        []string                       `tfsdk:"tags"`
	Variations  []featureResourceDataVariation `tfsdk:"variations"`
	Variables   []featureResourceDataVariable  `tfsdk:"variables"`
}

func (t featureResourceData) variationToSDK() []devcyclem.FeatureVariationDto {
	var variations []devcyclem.FeatureVariationDto
	for _, variation := range t.Variations {
		variations = append(variations, devcyclem.FeatureVariationDto{
			Key:       variation.Key.Value,
			Name:      variation.Name.Value,
			Variables: variation.Variables,
		})
	}
	return variations
}

func (t featureResourceData) variablesToSDK() []devcyclem.CreateVariableDto {
	var variables []devcyclem.CreateVariableDto
	for _, variable := range t.Variables {

		nvar := devcyclem.CreateVariableDto{
			Name:        variable.Name.Value,
			Description: variable.Description.Value,
			Key:         variable.Key.Value,
			Feature:     variable.FeatureKey.Value,
			Type_:       variable.Type.Value,
		}

		switch variable.Type.Value {
		case "string":
			var x interface{} = variable.DefaultStringValue.Value
			nvar.DefaultValue = &x
			break
		case "json":
			var x interface{} = variable.DefaultJsonValue.Value
			nvar.DefaultValue = &x
			break
		case "boolean":
			var x interface{} = variable.DefaultBoolValue.Value
			nvar.DefaultValue = &x
			break
		case "number":
			var x interface{} = variable.DefaultNumberValue.Value
			nvar.DefaultValue = &x
			break
		}

		variables = append(variables, nvar)
	}
	return variables
}

func variableToTF(vars []devcyclem.Variable) []featureResourceDataVariable {
	var variables []featureResourceDataVariable
	for _, variable := range vars {
		nvar := featureResourceDataVariable{
			Name:        types.String{Value: variable.Name},
			Description: types.String{Value: variable.Description},
			Key:         types.String{Value: variable.Key},
			FeatureKey:  types.String{Value: variable.Feature},
			Type:        types.String{Value: variable.Type_},
		}

		switch variable.Type_ {
		case "string":
			nvar.DefaultStringValue = types.String{Value: (*variable.DefaultValue).(string)}
			break
		case "json":
			nvar.DefaultJsonValue = types.String{Value: (*variable.DefaultValue).(string)}
			break
		case "boolean":
			nvar.DefaultBoolValue = types.Bool{Value: (*variable.DefaultValue).(bool)}
			break
		case "number":
			f := (*variable.DefaultValue).(big.Float)
			nvar.DefaultNumberValue = types.Number{Value: &f}
		}
		variables = append(variables, nvar)
	}
	return variables
}

type featureResourceDataVariable struct {
	Id                 types.String `tfsdk:"id"`
	Key                types.String `tfsdk:"key"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	FeatureKey         types.String `tfsdk:"feature_key"`
	Type               types.String `tfsdk:"type"`
	DefaultStringValue types.String `tfsdk:"default_string_value"`
	DefaultJsonValue   types.String `tfsdk:"default_json_value"`
	DefaultBoolValue   types.Bool   `tfsdk:"default_bool_value"`
	DefaultNumberValue types.Number `tfsdk:"default_number_value"`
}

type featureResourceDataVariation struct {
	Id        types.String           `tfsdk:"id"`
	Key       types.String           `tfsdk:"key"`
	Name      types.String           `tfsdk:"name"`
	Variables map[string]interface{} `tfsdk:"variables"`
}

func variationToTF(variations []devcyclem.Variation) []featureResourceDataVariation {
	var ret []featureResourceDataVariation
	for _, variation := range variations {
		nvar := featureResourceDataVariation{
			Key:       types.String{Value: variation.Key},
			Name:      types.String{Value: variation.Name},
			Variables: variation.Variables,
		}
		ret = append(ret, nvar)
	}
	return ret
}

type featureResource struct {
	provider provider
}

func (r featureResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data featureResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	feature, httpResponse, err := r.provider.MgmtClient.FeaturesApi.FeaturesControllerCreate(ctx, devcyclem.CreateFeatureDto{
		Name:        data.Name.Value,
		Key:         data.Key.Value,
		Description: data.Description.Value,
		Variations:  data.variationToSDK(),
		Variables:   data.variablesToSDK(),
		Type_:       data.Type.Value,
		Tags:        data.Tags,
	}, data.ProjectId.Value)
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create feature, got error: %s", err))
		return
	}

	data.Id.Value = feature.Id
	data.Key.Value = feature.Key
	data.Name.Value = feature.Name
	data.Description.Value = feature.Description
	data.ProjectId.Value = feature.Project
	data.Source.Value = feature.Source
	data.Type.Value = feature.Type_
	data.Tags = feature.Tags
	data.Variations = variationToTF(feature.Variations)
	data.Variables = variableToTF(feature.Variables)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r featureResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data featureResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	feature, httpResponse, err := r.provider.MgmtClient.FeaturesApi.FeaturesControllerFindOne(ctx, data.Key.Value, data.ProjectId.Value)
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read feature, got error: %s", err))
		return
	}

	data.Id.Value = feature.Id
	data.Key.Value = feature.Key
	data.Name.Value = feature.Name
	data.Description.Value = feature.Description
	data.ProjectId.Value = feature.Project
	data.Source.Value = feature.Source
	data.Type.Value = feature.Type_
	data.Tags = feature.Tags
	data.Variations = variationToTF(feature.Variations)
	data.Variables = variableToTF(feature.Variables)
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r featureResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data featureResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	feature, httpResponse, err := r.provider.MgmtClient.FeaturesApi.FeaturesControllerUpdate(ctx, devcyclem.UpdateFeatureDto{
		Name:        data.Name.Value,
		Key:         data.Key.Value,
		Description: data.Description.Value,
		Variations:  data.variationToSDK(),
		Variables:   data.variablesToSDK(),
		Type_:       data.Type.Value,
		Tags:        data.Tags,
	}, data.Key.Value, data.ProjectId.Value)
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update feature, got error: %s", err))
		return
	}

	data.Id.Value = feature.Id
	data.Key.Value = feature.Key
	data.Name.Value = feature.Name
	data.Description.Value = feature.Description
	data.ProjectId.Value = feature.Project
	data.Source.Value = feature.Source
	data.Type.Value = feature.Type_
	data.Tags = feature.Tags
	data.Variations = variationToTF(feature.Variations)
	data.Variables = variableToTF(feature.Variables)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r featureResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data featureResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := r.provider.MgmtClient.FeaturesApi.FeaturesControllerRemove(ctx, data.Key.Value, data.ProjectId.Value)
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete feature, got error: %s", err))

		return
	}

	resp.State.RemoveResource(ctx)
}

func (r featureResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
