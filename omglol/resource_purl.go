package omglol

import (
	"context"
	"strings"
	"time"

	"github.com/ejstreet/omglol-client-go/omglol"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pURLResource{}
	_ resource.ResourceWithConfigure   = &pURLResource{}
	_ resource.ResourceWithImportState = &pURLResource{}
)

// NewPURLResource is a helper function to simplify the provider implementation.
func NewPURLResource() resource.Resource {
	return &pURLResource{}
}

// purlResource is the resource implementation.
type pURLResource struct {
	client *omglol.Client
}

// pURLResourceModel maps the resource schema data.
type pURLResourceModel struct {
	Name      types.String `tfsdk:"name"`
	Address   types.String `tfsdk:"address"`
	URL       types.String `tfsdk:"url"`
	Listed    types.Bool   `tfsdk:"listed"`
	Counter   types.Int64  `tfsdk:"counter"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	ID        types.String `tfsdk:"id"`
}

// Metadata returns the resource type name.
func (r *pURLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_purl"
}

// Schema defines the schema for the resource.
func (r *pURLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage omg.lol Persistent URLs.",
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Your omg.lol address to create the pURL for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the PURL. The name field is how you will access your designated URL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The URL to link to.",
			},
			"listed": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Set true to list on your `address`.url.lol page.",
			},
			"counter": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of time a PURL has been accessed.",
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *pURLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pURLResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	purl := omglol.NewPersistentURL(plan.Name.ValueString(), plan.URL.ValueString(), plan.Listed.ValueBool())

	// Create Persistent URL
	err := r.client.CreatePersistentURL(plan.Address.ValueString(), *purl)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Persistent URL",
			"Could not create persistent URL, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Counter = types.Int64Value(0)
	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))
	plan.ID = types.StringValue(plan.Address.ValueString() + "_" + plan.Name.ValueString())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *pURLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pURLResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed PURL from omg.lol
	purl, err := r.client.GetPersistentURL(state.Address.ValueString(), state.Name.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Persistent URL",
			"Could not read persistent URL: "+err.Error(),
		)
		return
	}

	// Overwrite purl with refreshed state
	state.ID = types.StringValue(state.Address.ValueString() + "_" + state.Name.ValueString())
	state.Name = types.StringValue(purl.Name)
	state.URL = types.StringValue(purl.URL)
	state.Counter = types.Int64Value(int64(*purl.Counter))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *pURLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan pURLResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	purl := omglol.NewPersistentURL(plan.Name.ValueString(), plan.URL.ValueString(), plan.Listed.ValueBool())

	// Set account settings
	err := r.client.CreatePersistentURL(plan.Address.ValueString(), *purl)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Persistent URL",
			"Could not update persistent URL, unexpected error: "+err.Error(),
		)
		return
	}

	// Get refreshed pURL from omg.lol
	purl, err = r.client.GetPersistentURL(plan.Address.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Persistent URL",
			"Could not read persistent URL: "+err.Error(),
		)
		return
	}

	plan.Counter = types.Int64Value(int64(*purl.Counter))
	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))
	plan.ID = types.StringValue(plan.Address.ValueString() + "_" + plan.Name.ValueString())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *pURLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pURLResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing PURL
	err := r.client.DeletePersistentURL(state.Address.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting PURL",
			"Could not delete persistent URL, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *pURLResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*omglol.Client)
}

func (r *pURLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state pURLResourceModel

	// Retrieve import ID and save to id attribute
	state.ID = types.StringValue(req.ID)

	// Split ID into name and address
	parts := strings.Split(req.ID, "_")
	state.Address = types.StringValue(parts[0])
	state.Name = types.StringValue(parts[1])

	// Get refreshed pURL from omg.lol
	purl, err := r.client.GetPersistentURL(state.Address.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Persistent URL",
			"Could not read persistent URL: "+err.Error(),
		)
		return
	}

	// Overwrite purl with refreshed state

	state.Name = types.StringValue(purl.Name)
	state.URL = types.StringValue(purl.URL)
	state.Counter = types.Int64Value(int64(*purl.Counter))

	// Set refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
