package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/oleksii-kalinin/terraform-provider-sonarr/pkg/sonarr"
)

// SeriesDataSource implements the data source for finding an existing series in Sonarr.
type SeriesDataSource struct {
	client *sonarr.Client
}

// SeriesDataSourceModel describes the data source data model.
type SeriesDataSourceModel struct {
	ID               types.Int32  `tfsdk:"id"`
	Title            types.String `tfsdk:"title"`
	Path             types.String `tfsdk:"path"`
	QualityProfileId types.Int32  `tfsdk:"quality_profile_id"`
	Monitored        types.Bool   `tfsdk:"monitored"`
	SeasonFolder     types.Bool   `tfsdk:"season_folder"`
	TvdbId           types.Int32  `tfsdk:"tvdb_id"`
}

func (s *SeriesDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_series"
}

func (s *SeriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Data source for finding a series in Sonarr by title",
		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Title of the series to find",
			},
			"id": schema.Int32Attribute{
				Computed:    true,
				Description: "ID of the series in Sonarr",
			},
			"path": schema.StringAttribute{
				Computed:    true,
				Description: "Root folder path of the series",
			},
			"quality_profile_id": schema.Int32Attribute{
				Computed:    true,
				Description: "Quality profile ID for the series",
			},
			"monitored": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the series is monitored",
			},
			"season_folder": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to use season folders",
			},
			"tvdb_id": schema.Int32Attribute{
				Computed:    true,
				Description: "TVDB ID of the series",
			},
		},
	}
}

func (s *SeriesDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (s *SeriesDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data SeriesDataSourceModel

	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	allSeries, err := s.client.GetAllSeries()
	if err != nil {
		response.Diagnostics.AddError("Client error", fmt.Sprintf("Unable to get series from Sonarr: %s", err.Error()))
		return
	}

	searchTitle := strings.ToLower(data.Title.ValueString())
	var found *sonarr.Series
	for i := range allSeries {
		if strings.ToLower(allSeries[i].Title) == searchTitle {
			found = &allSeries[i]
			break
		}
	}

	if found == nil {
		response.Diagnostics.AddError("Series not found", fmt.Sprintf("No series found with title: %s", data.Title.ValueString()))
		return
	}

	data.ID = types.Int32Value(found.Id)
	data.Title = types.StringValue(found.Title)
	data.Path = types.StringValue(found.RootFolderPath)
	data.QualityProfileId = types.Int32Value(found.QualityProfileId)
	data.Monitored = types.BoolValue(found.Monitored)
	data.SeasonFolder = types.BoolValue(found.SeasonFolder)
	data.TvdbId = types.Int32Value(found.TvdbID)

	diags = response.State.Set(ctx, &data)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
}

// NewSeriesDataSource creates a new instance of the series data source.
func NewSeriesDataSource() datasource.DataSource {
	return &SeriesDataSource{}
}
