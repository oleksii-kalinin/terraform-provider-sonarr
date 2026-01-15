package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/oleksii-kalinin/terraform-provider-sonarr/pkg/sonarr"
)

type SystemStatusDataSource struct {
	client *sonarr.Client
}

type SystemStatusDataSourceModel struct {
	AppName types.String `tfsdk:"app_name"`
	Version types.String `tfsdk:"version"`
	OsName  types.String `tfsdk:"os_name"`
}

func (s *SystemStatusDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_system_status"
}

func (s *SystemStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "System information data source",
		Attributes: map[string]schema.Attribute{
			"app_name": schema.StringAttribute{
				Computed:    true,
				Description: "App name of the instance (Sonarr)",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "Version of the Sonarr installation",
			},
			"os_name": schema.StringAttribute{
				Computed:    true,
				Description: "OS name of the sonarr installation",
			},
		},
	}
}

func (s *SystemStatusDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*sonarr.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configuration type",
			fmt.Sprintf("Expected *sonarr.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData))
		return
	}
	s.client = client
}

func (s *SystemStatusDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data SystemStatusDataSourceModel

	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	status, err := s.client.GetSystemStatus()
	if err != nil {
		response.Diagnostics.AddError("Client error", fmt.Sprintf("Unable to communicate with Sonarr: %s", err.Error()))
		return
	}
	data.OsName = types.StringValue(status.OsName)
	data.Version = types.StringValue(status.Version)
	data.AppName = types.StringValue(status.AppName)

	diags = response.State.Set(ctx, &data)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
}

func NewSystemStatusDataSource() datasource.DataSource {
	return &SystemStatusDataSource{}
}
