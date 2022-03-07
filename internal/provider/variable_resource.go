package provider

import (
	"context"
	"fmt"
	devcyclem "github.com/devcyclehq/go-mgmt-sdk"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type variableResourceType struct{}

func (t variableResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Variable name",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "Variable description",
				Required:            true,
				Type:                types.StringType,
			},
			"key": {
				MarkdownDescription: "Variable key",
				Required:            true,
				Type:                types.StringType,
			},
			"feature_id": {
				MarkdownDescription: "Feature that this variable is attached to",
				Required:            true,
				Type:                types.StringType,
			},
			"project_id": {
				MarkdownDescription: "Project id that this feature and variable is attached to",
				Required:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Variable datatype",
				Required:            true,
				Type:                types.StringType,
			},
			"default_value": {
				MarkdownDescription: "Variable value. Will be parsed based on the type setting.",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Variable ID",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t variableResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return variableResource{
		provider: provider,
	}, diags
}

type variableResourceData struct {
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Key          types.String `tfsdk:"key"`
	FeatureId    types.String `tfsdk:"feature_id"`
	ProjectId    types.String `tfsdk:"project_id"`
	Type         types.String `tfsdk:"type"`
	DefaultValue types.String `tfsdk:"default_value"`
	Id           types.String `tfsdk:"id"`
}

type variableResource struct {
	provider provider
}

func (r variableResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data variableResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	variable, httpResponse, err := r.provider.MgmtClient.VariablesApi.VariablesControllerCreate(ctx, devcyclem.CreateVariableDto{
		Name:        data.Name.Value,
		Description: data.Description.Value,
		Key:         data.Key.Value,
		Feature:     data.FeatureId.Value,
		Type_:       data.Type.Value,
	}, data.ProjectId.Value)
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

	switch variable.Type_ {
	case "String":
		data.DefaultValue = types.String{Value: fmt.Sprintf("%v", variable.DefaultValue)}
		break
	case "JSON":
		data.DefaultValue = types.String{Value: fmt.Sprintf("%v", variable.DefaultValue)}
		break
	case "Boolean":
		fmt.Println(variable.DefaultValue)
		out, _ := strconv.ParseBool(fmt.Sprintf("%v", variable.DefaultValue))
		data.DefaultValue = types.String{Value: strconv.FormatBool(out)}
		break
	case "Number":
		out, _ := strconv.ParseFloat(fmt.Sprintf("%v", variable.DefaultValue), 64)
		data.DefaultValue = types.String{Value: strconv.FormatFloat(out, 'f', -1, 64)}
		break
	}
	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r variableResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data variableResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	variable, httpResponse, err := r.provider.MgmtClient.VariablesApi.VariablesControllerFindOne(ctx, data.Key.Value, data.ProjectId.Value)
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

	switch variable.Type_ {
	case "String":
		data.DefaultValue = types.String{Value: fmt.Sprintf("%v", variable.DefaultValue)}
		break
	case "JSON":
		data.DefaultValue = types.String{Value: fmt.Sprintf("%v", variable.DefaultValue)}
		break
	case "Boolean":
		fmt.Println(variable.DefaultValue)
		out, _ := strconv.ParseBool(fmt.Sprintf("%v", variable.DefaultValue))
		data.DefaultValue = types.String{Value: strconv.FormatBool(out)}
		break
	case "Number":
		out, _ := strconv.ParseFloat(fmt.Sprintf("%v", variable.DefaultValue), 64)
		data.DefaultValue = types.String{Value: strconv.FormatFloat(out, 'f', -1, 64)}
		break
	}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r variableResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data variableResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	variable, httpResponse, err := r.provider.MgmtClient.VariablesApi.VariablesControllerUpdate(ctx, devcyclem.UpdateVariableDto{
		Name:        data.Name.Value,
		Description: data.Description.Value,
		Key:         data.Key.Value,
		Feature:     data.FeatureId.Value,
	}, data.Id.Value, data.ProjectId.Value)
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
	switch variable.Type_ {
	case "String":
		data.DefaultValue = types.String{Value: fmt.Sprintf("%v", variable.DefaultValue)}
		break
	case "JSON":
		data.DefaultValue = types.String{Value: fmt.Sprintf("%v", variable.DefaultValue)}
		break
	case "Boolean":
		fmt.Println(variable.DefaultValue)
		out, _ := strconv.ParseBool(fmt.Sprintf("%v", variable.DefaultValue))
		data.DefaultValue = types.String{Value: strconv.FormatBool(out)}
		break
	case "Number":
		out, _ := strconv.ParseFloat(fmt.Sprintf("%v", variable.DefaultValue), 64)
		data.DefaultValue = types.String{Value: strconv.FormatFloat(out, 'f', -1, 64)}
		break
	}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r variableResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data variableResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := r.provider.MgmtClient.VariablesApi.VariablesControllerRemove(ctx, data.Key.Value, data.ProjectId.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r variableResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
