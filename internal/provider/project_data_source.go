package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type projectDataSourceType struct{}

func (t projectDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]tfsdk.Attribute{
			"key": {
				MarkdownDescription: "Project key, usually the lowercase, kebab case name of the project",
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

func (t projectDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return projectDataSource{
		provider: provider,
	}, diags
}

type projectDataSourceData struct {
	Name         types.String `tfsdk:"name"`
	Key          types.String `tfsdk:"key"`
	Description  types.String `tfsdk:"description"`
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	CreatedBy    types.String `tfsdk:"createdBy"`
	CreatedAt    types.String `tfsdk:"createdAt"`
	UpdatedAt    types.String `tfsdk:"updatedAt"`
}

type projectDataSource struct {
	provider provider
}

func (d projectDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data projectDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	project, httpResponse, err := d.provider.MgmtClient.ProjectsApi.ProjectsControllerFindOne(ctx, data.Key.Value)
	if err != nil || httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
		return
	}

	data.Name = types.String{Value: project.Name}
	data.Key = types.String{Value: project.Key}
	data.Organization = types.String{Value: project.Organization}
	data.CreatedBy = types.String{Value: project.CreatedBy}
	data.CreatedAt = types.String{Value: project.CreatedAt.String()}
	data.UpdatedAt = types.String{Value: project.UpdatedAt.String()}
	data.Id = types.String{Value: project.Id}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
