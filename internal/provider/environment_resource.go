package provider

import (
	"context"
	devcyclem "github.com/devcyclehq/go-mgmt-sdk"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type environmentResourceType struct{}

func (t environmentResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"project_id": {
				MarkdownDescription: "Project id or key of the project to which the environment belongs",
				Required:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Environment Name",
				Required:            true,
				Type:                types.StringType,
			},
			"key": {
				MarkdownDescription: "Environment Key",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "Environment Description",
				Required:            true,
				Type:                types.StringType,
			},
			"color": {
				MarkdownDescription: "Environment Color in Hex with leading #",
				Required:            true,
				Type:                types.StringType,
			},
			"type": {
				MarkdownDescription: "Environment Type",
				Required:            true,
				Type:                types.StringType,
			},
			"settings": {
				MarkdownDescription: "Environment Settings",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"app_icon_uri": {
						MarkdownDescription: "Environment App Icon Uri",
						Required:            true,
						Type:                types.StringType,
					},
				}),
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Environment Id",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
			"sdk_keys": {
				Computed:            true,
				MarkdownDescription: "SDK Keys for the environment",
				Type:                types.ListType{ElemType: types.StringType},
			},
		},
	}, nil
}

func (t environmentResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return environmentResource{
		provider: provider,
	}, diags
}

type environmentResourceData struct {
	Id          types.String                    `tfsdk:"id"`
	Key         types.String                    `tfsdk:"key"`
	Name        types.String                    `tfsdk:"name"`
	Description types.String                    `tfsdk:"description"`
	Color       types.String                    `tfsdk:"color"`
	Type        types.String                    `tfsdk:"type"`
	Settings    environmentResourceDataSettings `tfsdk:"settings"`
	ProjectId   types.String                    `tfsdk:"project_id"`
	SDKKeys     []string                        `tfsdk:"sdk_keys"`
}

func sdkKeyConvert(keys []devcyclem.ApiKey) []string {
	var sdkKeys []string
	for _, sdkKey := range keys {
		sdkKeys = append(sdkKeys, sdkKey.Key)
	}
	return sdkKeys
}

type environmentResourceDataSettings struct {
	AppIconURI types.String `tfsdk:"app_icon_uri"`
}

func (s *environmentResourceDataSettings) toCreateSDK() *devcyclem.AllOfCreateEnvironmentDtoSettings {
	return &devcyclem.AllOfCreateEnvironmentDtoSettings{
		AppIconURI: s.AppIconURI.Value,
	}
}
func (s *environmentResourceDataSettings) toUpdateSDK() *devcyclem.AllOfUpdateEnvironmentDtoSettings {
	return &devcyclem.AllOfUpdateEnvironmentDtoSettings{
		AppIconURI: s.AppIconURI.Value,
	}
}

type environmentResource struct {
	provider provider
}

func (r environmentResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data environmentResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	environment, httpResponse, err := r.provider.MgmtClient.EnvironmentsApi.EnvironmentsControllerCreate(ctx, devcyclem.CreateEnvironmentDto{
		Name:        data.Name.Value,
		Key:         data.Key.Value,
		Description: data.Description.Value,
		Color:       data.Color.Value,
		Type_:       data.Type.Value,
		Settings:    data.Settings.toCreateSDK(),
	}, data.ProjectId.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Id = types.String{Value: environment.Id}
	data.Key = types.String{Value: environment.Key}
	data.Name = types.String{Value: environment.Name}
	data.Description = types.String{Value: environment.Description}
	data.Color = types.String{Value: environment.Color}
	data.Type = types.String{Value: environment.Type_}
	data.Settings = environmentResourceDataSettings{
		AppIconURI: types.String{Value: environment.Settings.AppIconURI},
	}
	data.ProjectId = types.String{Value: environment.Project}
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Mobile)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Server)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Client)...)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r environmentResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data environmentResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	environment, httpResponse, err := r.provider.MgmtClient.EnvironmentsApi.EnvironmentsControllerFindOne(ctx, data.Key.Value, data.ProjectId.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}
	data.Id = types.String{Value: environment.Id}
	data.Key = types.String{Value: environment.Key}
	data.Name = types.String{Value: environment.Name}
	data.Description = types.String{Value: environment.Description}
	data.Color = types.String{Value: environment.Color}
	data.Type = types.String{Value: environment.Type_}
	data.Settings = environmentResourceDataSettings{
		AppIconURI: types.String{Value: environment.Settings.AppIconURI},
	}
	data.ProjectId = types.String{Value: environment.Project}
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Mobile)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Server)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Client)...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r environmentResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data environmentResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	environment, httpResponse, err := r.provider.MgmtClient.EnvironmentsApi.EnvironmentsControllerUpdate(ctx, devcyclem.UpdateEnvironmentDto{
		Name:        data.Name.Value,
		Key:         data.Key.Value,
		Description: data.Description.Value,
		Color:       data.Color.Value,
		Type_:       data.Type.Value,
		Settings:    data.Settings.toUpdateSDK(),
	}, data.Key.Value, data.ProjectId.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}

	data.Id = types.String{Value: environment.Id}
	data.Key = types.String{Value: environment.Key}
	data.Name = types.String{Value: environment.Name}
	data.Description = types.String{Value: environment.Description}
	data.Color = types.String{Value: environment.Color}
	data.Type = types.String{Value: environment.Type_}
	data.Settings = environmentResourceDataSettings{
		AppIconURI: types.String{Value: environment.Settings.AppIconURI},
	}
	data.ProjectId = types.String{Value: environment.Project}
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Mobile)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Server)...)
	data.SDKKeys = append(data.SDKKeys, sdkKeyConvert(environment.SdkKeys.Client)...)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r environmentResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data environmentResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := r.provider.MgmtClient.EnvironmentsApi.EnvironmentsControllerRemove(ctx, data.Key.Value, data.ProjectId.Value)
	if ret := handleDevCycleHTTP(err, httpResponse, &resp.Diagnostics); ret {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r environmentResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
