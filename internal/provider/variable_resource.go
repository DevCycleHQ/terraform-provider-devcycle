package provider

import (
	"context"
	"fmt"
	devcyclem "github.com/devcyclehq/go-mgmt-sdk"

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
			"feature_key": {
				MarkdownDescription: "Feature that this variable is attached to",
				Required:            true,
				Type:                types.StringType,
			},
			"project_key": {
				MarkdownDescription: "Project key that this feature and variable is attached to",
				Required:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Variable datatype",
				Required:            true,
				Type:                types.StringType,
			},
			"stringvalue": {
				MarkdownDescription: "Variable value if the type is string",
				Optional:            true,
				Type:                types.StringType,
			},
			"jsonvalue": {
				MarkdownDescription: "Variable value if the type is json",
				Optional:            true,
				Type:                types.StringType,
			},
			"boolvalue": {
				MarkdownDescription: "Variable value if the type is boolean",
				Optional:            true,
				Type:                types.BoolType,
			},
			"numvalue": {
				MarkdownDescription: "Variable value if the type is number",
				Optional:            true,
				Type:                types.NumberType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Variable ID",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
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
	FeatureId    types.String `tfsdk:"featureId"`
	ProjectId    types.String `tfsdk:"projectId"`
	Type         types.String `tfsdk:"type"`
	DefaultValue *interface{} `tfsdk:"default_value"`
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
		Name:         data.Name.Value,
		Description:  data.Description.Value,
		Key:          data.Key.Value,
		Feature:      data.FeatureId.Value,
		Type_:        data.Type.Value,
		DefaultValue: data.DefaultValue,
	}, data.ProjectId.Value)
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create variable, got error: %s", err))
		return
	}
	data.Id.Value = variable.Id
	data.Key.Value = variable.Key
	data.Name.Value = variable.Name
	data.Description.Value = variable.Description
	data.FeatureId.Value = variable.Feature
	data.ProjectId.Value = variable.Project
	data.Type.Value = variable.Type_
	data.DefaultValue = variable.DefaultValue

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
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create variable, got error: %s", err))
		return
	}
	data.Id.Value = variable.Id
	data.Key.Value = variable.Key
	data.Name.Value = variable.Name
	data.Description.Value = variable.Description
	data.FeatureId.Value = variable.Feature
	data.ProjectId.Value = variable.Project
	data.Type.Value = variable.Type_
	data.DefaultValue = variable.DefaultValue

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

	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create variable, got error: %s", err))
		return
	}
	data.Id.Value = variable.Id
	data.Key.Value = variable.Key
	data.Name.Value = variable.Name
	data.Description.Value = variable.Description
	data.FeatureId.Value = variable.Feature
	data.ProjectId.Value = variable.Project
	data.Type.Value = variable.Type_
	data.DefaultValue = variable.DefaultValue

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
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create variable, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r variableResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
