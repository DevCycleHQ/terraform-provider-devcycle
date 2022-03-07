package provider

import (
	"context"
	"fmt"
	dvc_server "github.com/devcyclehq/go-server-sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
)

type evaluatedNumberVariableDataSourceType struct{}

func (t evaluatedNumberVariableDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Evaluated Variable data source.",

		Attributes: map[string]tfsdk.Attribute{
			"user": userDataSchema(),
			"value": {
				MarkdownDescription: "Value of the Variable",
				Computed:            true,
				Type:                types.NumberType,
			},
			"default_value": {
				MarkdownDescription: "Default value of the Variable",
				Required:            true,
				Type:                types.NumberType,
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

func (t evaluatedNumberVariableDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return evaluatedNumberVariableDataSource{
		provider: provider,
	}, diags
}

type evaluatedNumberVariableDataSourceData struct {
	Id           types.String                        `tfsdk:"id"`
	Value        types.Number                        `tfsdk:"value"`
	User         evaluatedVariableDataSourceDataUser `tfsdk:"user"`
	DefaultValue types.Number                        `tfsdk:"default_value"`
}

type evaluatedNumberVariableDataSource struct {
	provider provider
}

func (d evaluatedNumberVariableDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data evaluatedNumberVariableDataSourceData

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
	data.Value = types.Number{Value: big.NewFloat((*variable.Value).(float64))}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
