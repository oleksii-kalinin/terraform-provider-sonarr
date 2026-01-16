package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/oleksii-kalinin/terraform-provider-sonarr/pkg/sonarr"
)

type SeriesResource struct {
	client *sonarr.Client
}

type SeriesResourceModel struct {
	ID               types.String     `tfsdk:"id"`
	TvdbId           types.Int32      `tfsdk:"tvdb_id"`
	Title            types.String     `tfsdk:"title"`
	Path             types.String     `tfsdk:"path"`
	Monitored        types.Bool       `tfsdk:"monitored"`
	QualityProfileId types.Int32      `tfsdk:"quality_profile"`
	AddOptions       *AddOptionsModel `tfsdk:"add_options"`
}

type AddOptionsModel struct {
	Monitor types.String `tfsdk:"monitor"`
}

func (s *SeriesResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_series"
}

func (s *SeriesResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Resource for the Sonarr Series",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tvdb_id": schema.Int32Attribute{
				Required:      true,
				PlanModifiers: []planmodifier.Int32{int32planmodifier.RequiresReplace()},
			},
			"title": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"path": schema.StringAttribute{
				Required: true,
			},
			"monitored": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"quality_profile": schema.Int32Attribute{
				Required: true,
			},
			"add_options": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"monitor": schema.StringAttribute{
						Optional:    true,
						Description: "Valid values: all, future, missing, etc.",
					},
				},
			},
		},
	}
}

func (s *SeriesResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan SeriesResourceModel
	diags := request.Plan.Get(ctx, &plan)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	var addOpts *sonarr.AddOptions
	if plan.AddOptions != nil {
		addOpts = &sonarr.AddOptions{
			Monitor: plan.AddOptions.Monitor.ValueString(),
		}
	}

	seriesReq := sonarr.Series{
		Title:            plan.Title.ValueString(),
		TvdbID:           plan.TvdbId.ValueInt32(),
		QualityProfileId: plan.QualityProfileId.ValueInt32(),
		RootFolderPath:   plan.Path.ValueString(),
		Monitored:        plan.Monitored.ValueBool(),
		AddOptions:       addOpts,
	}

	seriesRes, err := s.client.CreateSeries(&seriesReq)
	if err != nil {
		response.Diagnostics.AddError("Error creating series", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(int(seriesRes.Id)))

	diags = response.State.Set(ctx, &plan)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
}

func (s *SeriesResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state SeriesResourceModel
	diags := request.State.Get(ctx, &state)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		response.State.RemoveResource(ctx)
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing series ID", err.Error())
		return
	}

	seriesReq, err := s.client.GetSeries(id)
	if err != nil {
		response.Diagnostics.AddError("Error getting series", err.Error())
		return
	}

	if seriesReq == nil {
		response.State.RemoveResource(ctx)
		return
	}

	state.TvdbId = types.Int32Value(seriesReq.TvdbID)
	state.Title = types.StringValue(seriesReq.Title)
	state.Path = types.StringValue(seriesReq.RootFolderPath)
	state.Monitored = types.BoolValue(seriesReq.Monitored)
	state.QualityProfileId = types.Int32Value(seriesReq.QualityProfileId)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (s *SeriesResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state SeriesResourceModel

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Error parsing series ID from the state", err.Error())
		return
	}

	currentSeries, err := s.client.GetSeries(id)
	if err != nil {
		response.Diagnostics.AddError("Error fetching series", err.Error())
		return
	}
	if currentSeries == nil {
		response.Diagnostics.AddError("Series not found", "Could not find series to update. It might have been deleted manually.")
		return
	}

	currentSeries.Title = plan.Title.ValueString()
	currentSeries.Monitored = plan.Monitored.ValueBool()
	currentSeries.RootFolderPath = plan.Path.ValueString()
	currentSeries.QualityProfileId = plan.QualityProfileId.ValueInt32()
	currentSeries.TvdbID = plan.TvdbId.ValueInt32()

	_, err = s.client.UpdateSeries(currentSeries)
	if err != nil {
		response.Diagnostics.AddError("Error updating series", err.Error())
		return
	}

	plan.ID = state.ID
	response.State.Set(ctx, &plan)
}

func (s *SeriesResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state SeriesResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() || state.ID.ValueString() == "" {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Invalid ID format", err.Error())
		return
	}

	tflog.Info(ctx, "Deleting series", map[string]any{"id": id, "title": state.Title.ValueString()})

	err = s.client.DeleteSeries(id, true)
	if err != nil {
		response.Diagnostics.AddError("Error Deleting Series", err.Error())
		return
	}
}

func (s *SeriesResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*sonarr.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sonarr.Client, got: %T", request.ProviderData),
		)
		return
	}

	s.client = client
}

func NewSeriesResource() resource.Resource {
	return &SeriesResource{}
}
