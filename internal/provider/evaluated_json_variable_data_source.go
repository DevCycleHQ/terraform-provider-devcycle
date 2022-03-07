package provider

import (
	"context"
	"encoding/json"
	"fmt"
	dvc_server "github.com/devcyclehq/go-server-sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type evaluatedJSONVariableDataSourceType struct{}

func (t evaluatedJSONVariableDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Evaluated Variable data source.",

		Attributes: map[string]tfsdk.Attribute{
			"user": userDataSchema(),
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

func (t evaluatedJSONVariableDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return evaluatedJSONVariableDataSource{
		provider: provider,
	}, diags
}

type evaluatedJSONVariableDataSourceData struct {
	Id           types.String                        `tfsdk:"id"`
	Value        types.String                        `tfsdk:"value"`
	User         evaluatedVariableDataSourceDataUser `tfsdk:"user"`
	DefaultValue types.String                        `tfsdk:"default_value"`
}

type evaluatedJSONVariableDataSource struct {
	provider provider
}

func (d evaluatedJSONVariableDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data evaluatedJSONVariableDataSourceData

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

	jsonstring, err := json.Marshal(*variable.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Variable, got error: %s", err))
		return
	}

	data.Id = types.String{Value: variable.Id}
	data.Value = types.String{Value: string(jsonstring)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
