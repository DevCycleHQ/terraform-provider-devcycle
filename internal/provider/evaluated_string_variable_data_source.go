package provider

import (
	"context"
	"fmt"
	dvc_server "github.com/devcyclehq/go-server-sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type evaluatedStringVariableDataSourceType struct{}

func (t evaluatedStringVariableDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Evaluated Variable data source.",

		Attributes: map[string]tfsdk.Attribute{
			"user": {
				MarkdownDescription: "User data to drive bucketing into variations",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						MarkdownDescription: "User ID",
						Required:            true,
						Type:                types.StringType,
					},
					"name": {
						MarkdownDescription: "User name",
						Optional:            true,
						Type:                types.StringType,
					},
					"app_version": {
						MarkdownDescription: "User app version",
						Optional:            true,
						Type:                types.StringType,
					},
					"email": {
						MarkdownDescription: "User email",
						Optional:            true,
						Type:                types.StringType,
					},
					"app_build": {
						MarkdownDescription: "User app build",
						Optional:            true,
						Type:                types.StringType,
					},
				}),
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"value": {
				MarkdownDescription: "Value of the Variable",
				Computed:            true,
				Type:                types.StringType,
			},
			"default_value": {
				MarkdownDescription: "Default value of the Variable",
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Required:            true,
				MarkdownDescription: "Variable ID",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t evaluatedStringVariableDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return evaluatedStringVariableDataSource{
		provider: provider,
	}, diags
}

type evaluatedStringVariableDataSourceData struct {
	Id           types.String                        `tfsdk:"id"`
	Value        types.String                        `tfsdk:"value"`
	User         evaluatedVariableDataSourceDataUser `tfsdk:"user"`
	DefaultValue types.String                        `tfsdk:"default_value"`
}

type evaluatedStringVariableDataSource struct {
	provider provider
}

func (d evaluatedStringVariableDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data evaluatedStringVariableDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	userData := dvc_server.UserData{
		UserId: "" + data.User.Id.Value,
	}

	variable, err := d.provider.ServerClient.DevcycleApi.Variable(d.provider.ServerClientContext, userData, data.Id.Value, data.DefaultValue.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Variable, got error: %s", err))
		return
	}

	data.Id = types.String{Value: variable.Id}
	data.Value = types.String{Value: (*variable.Value).(string)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
