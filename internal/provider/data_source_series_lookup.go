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

// SeriesLookupDataSource implements the data source for searching series on TVDB via Sonarr.
// This is used to find series information before adding them to the Sonarr library.
type SeriesLookupDataSource struct {
	client *sonarr.Client
}

// SeriesLookupDataSourceModel describes the data source data model.
type SeriesLookupDataSourceModel struct {
	Term        types.String `tfsdk:"term"`
	Title       types.String `tfsdk:"title"`
	SortTitle   types.String `tfsdk:"sort_title"`
	Status      types.String `tfsdk:"status"`
	Overview    types.String `tfsdk:"overview"`
	Network     types.String `tfsdk:"network"`
	Year        types.Int32  `tfsdk:"year"`
	TvdbId      types.Int32  `tfsdk:"tvdb_id"`
	ImdbId      types.String `tfsdk:"imdb_id"`
	Runtime     types.Int32  `tfsdk:"runtime"`
	SeasonCount types.Int32  `tfsdk:"season_count"`
}

func (s *SeriesLookupDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_series_lookup"
}

func (s *SeriesLookupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Data source for looking up series information from TVDB via Sonarr. Use this to find series details before adding them.",
		Attributes: map[string]schema.Attribute{
			"term": schema.StringAttribute{
				Required:    true,
				Description: "Search term to find the series (searches TVDB)",
			},
			"title": schema.StringAttribute{
				Computed:    true,
				Description: "Title of the series",
			},
			"sort_title": schema.StringAttribute{
				Computed:    true,
				Description: "Sort title of the series",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the series (continuing, ended, etc.)",
			},
			"overview": schema.StringAttribute{
				Computed:    true,
				Description: "Overview/description of the series",
			},
			"network": schema.StringAttribute{
				Computed:    true,
				Description: "Network the series airs on",
			},
			"year": schema.Int32Attribute{
				Computed:    true,
				Description: "Year the series started",
			},
			"tvdb_id": schema.Int32Attribute{
				Computed:    true,
				Description: "TVDB ID of the series",
			},
			"imdb_id": schema.StringAttribute{
				Computed:    true,
				Description: "IMDB ID of the series",
			},
			"runtime": schema.Int32Attribute{
				Computed:    true,
				Description: "Runtime of episodes in minutes",
			},
			"season_count": schema.Int32Attribute{
				Computed:    true,
				Description: "Number of seasons",
			},
		},
	}
}

func (s *SeriesLookupDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (s *SeriesLookupDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data SeriesLookupDataSourceModel

	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	results, err := s.client.LookupSeries(data.Term.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Client error", fmt.Sprintf("Unable to lookup series: %s", err.Error()))
		return
	}

	if len(results) == 0 {
		response.Diagnostics.AddError("Series not found", fmt.Sprintf("No series found matching: %s", data.Term.ValueString()))
		return
	}

	// Find exact match first, otherwise use first result
	searchTerm := strings.ToLower(data.Term.ValueString())
	var found *sonarr.SeriesLookup
	for i := range results {
		if strings.ToLower(results[i].Title) == searchTerm {
			found = &results[i]
			break
		}
	}
	if found == nil {
		found = &results[0]
	}

	data.Title = types.StringValue(found.Title)
	data.SortTitle = types.StringValue(found.SortTitle)
	data.Status = types.StringValue(found.Status)
	data.Overview = types.StringValue(found.Overview)
	data.Network = types.StringValue(found.Network)
	data.Year = types.Int32Value(found.Year)
	data.TvdbId = types.Int32Value(found.TvdbId)
	data.ImdbId = types.StringValue(found.ImdbId)
	data.Runtime = types.Int32Value(found.Runtime)
	data.SeasonCount = types.Int32Value(found.SeasonCount)

	diags = response.State.Set(ctx, &data)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
}

// NewSeriesLookupDataSource creates a new instance of the series lookup data source.
func NewSeriesLookupDataSource() datasource.DataSource {
	return &SeriesLookupDataSource{}
}
