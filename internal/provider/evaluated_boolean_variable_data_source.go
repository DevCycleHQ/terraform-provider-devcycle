package provider

import (
	"context"
	"fmt"
	dvc_server "github.com/devcyclehq/go-server-sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type evaluatedBoolVariableDataSourceType struct{}

func (t evaluatedBoolVariableDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Evaluated Variable data source.",

		Attributes: map[string]tfsdk.Attribute{
			"user": userDataSchema(),
			"value": {
				MarkdownDescription: "Value of the Variable. Either true or false.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"default_value": {
				MarkdownDescription: "Default value of the Variable. Either true or false.",
				Required:            true,
				Type:                types.BoolType,
			},
			"id": {
				Required:            true,
				MarkdownDescription: "Variable ID or key. Recommended to use the key when not managing an entire project in Terraform.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t evaluatedBoolVariableDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return evaluatedBooleanVariableDataSource{
		provider: provider,
	}, diags
}

type evaluatedBooleanVariableDataSourceData struct {
	Id           types.String                        `tfsdk:"id"`
	Value        types.Bool                          `tfsdk:"value"`
	User         evaluatedVariableDataSourceDataUser `tfsdk:"user"`
	DefaultValue types.Bool                          `tfsdk:"default_value"`
}

type evaluatedBooleanVariableDataSource struct {
	provider provider
}

func (d evaluatedBooleanVariableDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data evaluatedBooleanVariableDataSourceData

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
	data.Value = types.Bool{Value: (*variable.Value).(bool)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
