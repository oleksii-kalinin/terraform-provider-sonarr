package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/oleksii-kalinin/terraform-provider-sonarr/pkg/sonarr"
)

type SonarrProvider struct {
	version string
}

type SonarrProviderModel struct {
	Url    types.String `tfsdk:"url"`
	ApiKey types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SonarrProvider{
			version: version,
		}
	}
}

func (sp *SonarrProvider) Metadata(_ context.Context, _ provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "sonarr"
	res.Version = sp.version
}

func (sp *SonarrProvider) Schema(_ context.Context, _ provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Interact with Sonarr via Terraform",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "URL of the Sonarr server. Can also be set via SONARR_URL environment variable.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "API key for the sonarr instance. Can also be set via SONARR_API_KEY environment variable.",
				Optional:    true,
			},
		},
	}
}

func (sp *SonarrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	var config SonarrProviderModel
	diag := req.Config.Get(ctx, &config)

	res.Diagnostics.Append(diag...)
	if res.Diagnostics.HasError() {
		return
	}

	if config.Url.IsUnknown() || config.ApiKey.IsUnknown() {
		return
	}

	if config.Url.IsNull() {
		if v := os.Getenv("SONARR_URL"); v != "" {
			config.Url = types.StringValue(v)
		} else {
			res.Diagnostics.AddError("Sonarr URL missing", "Sonarr URL should be provided")
			return
		}
	}

	if config.ApiKey.IsNull() {
		if v := os.Getenv("SONARR_API_KEY"); v != "" {
			config.ApiKey = types.StringValue(v)
		} else {
			res.Diagnostics.AddError("Sonarr API key missing", "Sonarr API key should be provided")
			return
		}
	}

	client := sonarr.NewClient(config.Url.ValueString(), config.ApiKey.ValueString())

	res.DataSourceData = client
	res.ResourceData = client
}

func (sp *SonarrProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (sp *SonarrProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
