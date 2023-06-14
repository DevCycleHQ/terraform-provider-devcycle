package provider

import (
	"context"
	devcyclem "github.com/devcyclehq/go-mgmt-sdk"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type projectResourceType struct{}

func (t projectResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DevCycle project resource. Allows for creation/modification of a project.",

		Attributes: map[string]tfsdk.Attribute{
			"description": {
				MarkdownDescription: "Description of the project",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Name of the project",
				Required:            true,
				Type:                types.StringType,
			},
			"key": {
				MarkdownDescription: "Project key, usually the lowercase, kebab case name of the project",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Project Key",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
			"organization": {
				MarkdownDescription: "Organization that the project belongs to",
				Computed:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (t projectResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return projectResource{
		provider: provider,
	}, diags
}

type projectResourceData struct {
	Name         types.String `tfsdk:"name"`
	Key          types.String `tfsdk:"key"`
	Description  types.String `tfsdk:"description"`
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
}

type projectResource struct {
	provider provider
}

func (r projectResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data projectResourceData

	if !r.provider.configured {
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

	project, httpResponse, err := r.provider.MgmtClient.ProjectsApi.ProjectsControllerCreate(ctx, devcyclem.CreateProjectDto{
		Name:        data.Name.Value,
		Key:         strings.ToLower(data.Key.Value),
		Description: data.Description.Value,
	})
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Name = types.String{Value: project.Name}
	data.Key = types.String{Value: project.Key}
	data.Organization = types.String{Value: project.Organization}
	data.Id = types.String{Value: project.Id}

	tflog.Trace(ctx, "Created a project with id %s", project.Id)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r projectResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data projectResourceData
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

	project, httpResponse, err := r.provider.MgmtClient.ProjectsApi.ProjectsControllerFindOne(ctx, data.Key.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Name = types.String{Value: project.Name}
	data.Key = types.String{Value: project.Key}
	data.Organization = types.String{Value: project.Organization}
	data.Id = types.String{Value: project.Id}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r projectResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data projectResourceData
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

	project, httpResponse, err := r.provider.MgmtClient.ProjectsApi.ProjectsControllerUpdate(ctx, devcyclem.UpdateProjectDto{
		Name:        data.Name.Value,
		Key:         data.Key.Value,
		Description: data.Description.Value,
	}, data.Key.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Name = types.String{Value: project.Name}
	data.Key = types.String{Value: project.Key}
	data.Organization = types.String{Value: project.Organization}
	data.Id = types.String{Value: project.Id}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r projectResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data projectResourceData
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

	httpResponse, err := r.provider.MgmtClient.ProjectsApi.ProjectsControllerRemove(ctx, data.Key.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	tflog.Trace(ctx, "Deleted project with id %s", data.Key.Value)

	resp.State.RemoveResource(ctx)
}

func (r projectResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
