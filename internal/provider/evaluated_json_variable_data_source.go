package provider

import (
	"context"
	"encoding/json"
	"fmt"
	dvc_server "github.com/devcyclehq/go-server-sdk/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type evaluatedJSONVariableDataSourceType struct{}

func (t evaluatedJSONVariableDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Evaluated Variable data source. Each instance of this data source represents a single evaluated variable, under a single userdata context.",

		Attributes: map[string]tfsdk.Attribute{
			"user": userDataSchema(),
			"value": {
				MarkdownDescription: "Value of the Variable",
				Computed:            true,
				Type:                types.StringType,
			},
			"default_value": {
				MarkdownDescription: "Default value of the Variable. Used as a fallback in case there is no variation value set.",
				Required:            true,
				Type:                types.StringType,
			},
			"key": {
				Required:            true,
				MarkdownDescription: "Variable ID or key. Recommended to use the key when not managing an entire project in Terraform.",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
			"id": {
				Computed: true,
				Type:     types.StringType,
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
	Key          types.String                        `tfsdk:"key"`
	Value        types.String                        `tfsdk:"value"`
	User         evaluatedVariableDataSourceDataUser `tfsdk:"user"`
	DefaultValue types.String                        `tfsdk:"default_value"`
	Id           types.String                        `tfsdk:"id"`
}

type evaluatedJSONVariableDataSource struct {
	provider provider
}

func (d evaluatedJSONVariableDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data evaluatedJSONVariableDataSourceData
	if !d.provider.configured {
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

	userData := dvc_server.DVCUser{
		UserId: "" + data.User.Id.Value,
	}

	defaultValueJSON := []byte(data.DefaultValue.Value)
	var defaultValue map[string]any
	err := json.Unmarshal(defaultValueJSON, &defaultValue)
	if err != nil {
		resp.Diagnostics.AddError("JSON Serialization Error", fmt.Sprintf("Unable to read Variable, got error: %s", err))
		return
	}
	variable, err := d.provider.ServerClient.Variable(userData, data.Key.Value, defaultValue)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Variable, got error: %s", err))
		return
	}

	jsonstring, err := json.Marshal(variable.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Variable, got error: %s", err))
		return
	}

	data.Key = types.String{Value: variable.Key}
	data.Id = data.Key

	data.Value = types.String{Value: string(jsonstring)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
