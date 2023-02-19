package omglol

import (
	"context"
	"time"

	"github.com/ejstreet/omglol-client-go/omglol"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &accountSettingsResource{}
	_ resource.ResourceWithConfigure = &accountSettingsResource{}
)

// NewAccountSettingsResource is a helper function to simplify the provider implementation.
func NewAccountSettingsResource() resource.Resource {
	return &accountSettingsResource{}
}

// accountsettingsResource is the resource implementation.
type accountSettingsResource struct {
	client *omglol.Client
}

// accountSettingsResourceModel maps the resource schema data.
type accountSettingsResourceModel struct {
	Communication types.String `tfsdk:"communication"`
	DateFormat    types.String `tfsdk:"date_format"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	ID            types.String `tfsdk:"id"`
}

// Metadata returns the resource type name.
func (r *accountSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_settings"
}

// Schema defines the schema for the resource.
func (r *accountSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"communication": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Commuinication preferences. Valid values are `email_ok` and `email_not_ok`",
				Validators: []validator.String{
					stringvalidator.OneOf("email_ok", "email_not_ok"),
				},
			},
			"date_format": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Date preferences. Valid values are: `iso_8601` for *YYYY-MM-DD*, `dmy` for *DD-MM-YYYY*, and `mdy` for *MM-DD-YYYY*.",
				Validators: []validator.String{
					stringvalidator.OneOf("iso_8601", "dmy", "mdy"),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *accountSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan accountSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan, and compute owner data
	settings := map[string]string{
		"communication": plan.Communication.ValueString(),
		"date_format":   plan.DateFormat.ValueString(),
	}

	// Set account settings
	err := r.client.SetAccountSettings(settings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating settings",
			"Could not update settings, unexpected error: "+err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.ID = types.StringValue("_")

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *accountSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state accountSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed account settings from omg.lol
	settings, err := r.client.GetAccountSettings()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Account Settings",
			"Could not read Account Settings: "+err.Error(),
		)
		return
	}

	// Overwrite settings with refreshed state
	state.Communication = types.StringValue(settings.Communication)
	state.DateFormat = types.StringValue(settings.DateFormat)
	state.ID = types.StringValue("_")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *accountSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan accountSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan, and compute owner data
	settings := map[string]string{
		"communication": plan.Communication.ValueString(),
		"date_format":   plan.DateFormat.ValueString(),
	}

	// Set account settings
	err := r.client.SetAccountSettings(settings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating settings",
			"Could not update settings, unexpected error: "+err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.ID = types.StringValue("_")

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *accountSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Configure adds the provider configured client to the resource.
func (r *accountSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*omglol.Client)
}
