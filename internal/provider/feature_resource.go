package provider

import (
	"context"
	"encoding/json"
	"fmt"
	devcyclem "github.com/devcyclehq/go-mgmt-sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"sort"
	"strconv"
)

type featureResourceType struct{}

func (t featureResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DevCycle Feature resource. It's recommended to use the variable resource instead of this resource to manage variables. Variations currently have to be managed in this resource.",

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
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"project_id": {
				MarkdownDescription: "Project ID that the feature belongs to",
				Required:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Feature Type",
				Required:            true,
				Type:                types.StringType,
			},
			"source": {
				MarkdownDescription: "Source of Feature creation",
				Computed:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "Feature tags",
				Optional:            true,
				Type:                types.ListType{ElemType: types.StringType},
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
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
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
	ProjectId   types.String                   `tfsdk:"project_id"`
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
			Variables: variation.variationMapTypeFix(t.Variables),
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
			Feature:     t.Key.Value,
			Type_:       variable.Type.Value,
		}
		variables = append(variables, nvar)
	}
	return variables
}

type featureResourceDataVariable struct {
	Id          types.String `tfsdk:"id"`
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	FeatureKey  types.String `tfsdk:"feature_key"`
	Type        types.String `tfsdk:"type"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

type featureResourceDataVariation struct {
	Id        types.String      `tfsdk:"id"`
	Key       types.String      `tfsdk:"key"`
	Name      types.String      `tfsdk:"name"`
	Variables map[string]string `tfsdk:"variables"`
}

func variationMapToString(variations map[string]interface{}, variables []featureResourceDataVariable) map[string]string {
	ret := make(map[string]string)
	for k, v := range variations {
		for _, variable := range variables {
			if variable.Key.Value == k {
				switch variable.Type.Value {
				case "String":
					ret[k] = v.(string)
					break
				case "Number":
					ret[k] = fmt.Sprintf("%f", v.(float64))
					break
				case "Boolean":
					ret[k] = fmt.Sprintf("%t", v.(bool))
					break
				case "JSON":
					marshal, err := json.Marshal(v)
					if err != nil {
						tflog.Error(context.Background(), "Error parsing json", err)
						return nil
					}
					ret[k] = string(marshal)
					break
				}
			}
		}
		ret[k] = fmt.Sprintf("%v", v)
	}
	return ret
}

func (f *featureResourceDataVariation) variationMapTypeFix(variables []featureResourceDataVariable) map[string]interface{} {
	ret := make(map[string]interface{})

	for k, v := range f.Variables {
		for _, variable := range variables {
			if variable.Key.Value == k {
				switch variable.Type.Value {
				case "String":
					ret[k] = v
					break
				case "JSON":
					var marshalled interface{}
					err := json.Unmarshal([]byte(v), &marshalled)
					if err != nil {
						tflog.Error(context.Background(), "Error parsing json", err)
						return nil
					}
					break
				case "Number":
					float, err := strconv.ParseFloat(v, 64)
					if err != nil {
						tflog.Error(context.Background(), "Error parsing float", err)
						return nil
					}
					ret[k] = float
					break
				case "Boolean":
					b, err := strconv.ParseBool(v)
					if err != nil {
						tflog.Error(context.Background(), "Error parsing bool", err)
						return nil
					}
					ret[k] = b
					break
				}
			}
		}
	}
	return ret
}

func variationToTF(variations []devcyclem.Variation, variables []featureResourceDataVariable) []featureResourceDataVariation {
	var ret []featureResourceDataVariation
	for _, variation := range variations {
		nvar := featureResourceDataVariation{
			Key:       types.String{Value: variation.Key},
			Name:      types.String{Value: variation.Name},
			Variables: variationMapToString(variation.Variables, variables),
			Id:        types.String{Value: variation.Id},
		}
		ret = append(ret, nvar)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Key.Value < ret[j].Key.Value
	})
	return ret
}

func variableToTF(vars []devcyclem.Variable) []featureResourceDataVariable {
	var ret []featureResourceDataVariable
	for _, variable := range vars {
		nvar := featureResourceDataVariable{
			Key:         types.String{Value: variable.Key},
			Name:        types.String{Value: variable.Name},
			Description: types.String{Value: variable.Description},
			FeatureKey:  types.String{Value: variable.Feature},
			Type:        types.String{Value: variable.Type_},
			Id:          types.String{Value: variable.Id},
			CreatedAt:   types.String{Value: variable.CreatedAt.String()},
			UpdatedAt:   types.String{Value: variable.UpdatedAt.String()},
		}
		ret = append(ret, nvar)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Key.Value < ret[j].Key.Value
	})
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
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. Authentication is required to be configured.",
		)
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
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Id = types.String{Value: feature.Id}
	data.Key = types.String{Value: feature.Key}
	data.Name = types.String{Value: feature.Name}
	data.Description = types.String{Value: feature.Description}
	data.Type = types.String{Value: feature.Type_}
	data.Tags = feature.Tags
	data.ProjectId = types.String{Value: feature.Project}
	data.Source = types.String{Value: feature.Source}
	data.Variables = variableToTF(feature.Variables)
	data.Variations = variationToTF(feature.Variations, data.Variables)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r featureResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data featureResourceData
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. Authentication is required to be configured.",
		)
		return
	}
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	feature, httpResponse, err := r.provider.MgmtClient.FeaturesApi.FeaturesControllerFindOne(ctx, data.Key.Value, data.ProjectId.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Id = types.String{Value: feature.Id}
	data.Key = types.String{Value: feature.Key}
	data.Name = types.String{Value: feature.Name}
	data.Description = types.String{Value: feature.Description}
	data.Type = types.String{Value: feature.Type_}
	data.Tags = feature.Tags
	data.ProjectId = types.String{Value: feature.Project}
	data.Source = types.String{Value: feature.Source}
	data.Variables = variableToTF(feature.Variables)
	data.Variations = variationToTF(feature.Variations, data.Variables)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r featureResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data featureResourceData
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. Authentication is required to be configured.",
		)
		return
	}
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	feature, httpResponse, err := r.provider.MgmtClient.FeaturesApi.FeaturesControllerUpdate(ctx, devcyclem.UpdateFeatureDto{
		Name:        data.Name.Value,
		Key:         data.Key.Value,
		Description: data.Description.Value,
		Type_:       data.Type.Value,
		Tags:        data.Tags,
		Variables:   data.variablesToSDK(),
		Variations:  data.variationToSDK(),
	}, data.Key.Value, data.ProjectId.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Id = types.String{Value: feature.Id}
	data.Key = types.String{Value: feature.Key}
	data.Name = types.String{Value: feature.Name}
	data.Description = types.String{Value: feature.Description}
	data.Type = types.String{Value: feature.Type_}
	data.Tags = feature.Tags
	data.ProjectId = types.String{Value: feature.Project}
	data.Source = types.String{Value: feature.Source}
	data.Variables = variableToTF(feature.Variables)
	data.Variations = variationToTF(feature.Variations, data.Variables)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r featureResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data featureResourceData
	if !r.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. Authentication is required to be configured.",
		)
		return
	}
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	for _, variable := range data.Variables {
		httpResponse, err := r.provider.MgmtClient.VariablesApi.VariablesControllerRemove(ctx, variable.Id.Value, data.ProjectId.Value)
		if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
			return
		}
	}
	httpResponse, err := r.provider.MgmtClient.FeaturesApi.FeaturesControllerRemove(ctx, data.Key.Value, data.ProjectId.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r featureResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
