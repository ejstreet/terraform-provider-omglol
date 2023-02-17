package omglol

import (
	"context"
	"strings"

	"github.com/ejstreet/omglol-client-go/omglol"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dnsRecordResource{}
	_ resource.ResourceWithConfigure   = &dnsRecordResource{}
	_ resource.ResourceWithImportState = &dnsRecordResource{}
)

// NewDNSRecordResource is a helper function to simplify the provider implementation.
func NewDNSRecordResource() resource.Resource {
	return &dnsRecordResource{}
}

// dnsrecordResource is the resource implementation.
type dnsRecordResource struct {
	client *omglol.Client
}

// orderResourceModel maps the resource schema data.
type dnsRecordResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Address   types.String `tfsdk:"address"`
	Name      types.String `tfsdk:"name"`
	Data      types.String `tfsdk:"data"`
	Priority  types.Int64  `tfsdk:"priority"`
	TTL       types.Int64  `tfsdk:"ttl"`
	FQDN      types.String `tfsdk:"fqdn"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *dnsRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

// Schema defines the schema for the resource.
func (r *dnsRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage omg.lol DNS records.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The record type. Valid values are `A`, `AAAA`, `CAA`, `CNAME`, `TXT`, `MX`, `NS`, and `SRV`.",
				Validators: []validator.String{
					stringvalidator.OneOf("A", "AAAA", "CAA", "CNAME", "TXT", "MX", "NS", "SRV"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Your omg.lol address to create the record for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The prefix to attach before the address. Enter `@` to use the top level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The data to enter into the record.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ttl": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The Time-To-Live (TTL) of the record.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The priority of the record. Only applies to MX records.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"fqdn": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The fully qualified domain name of the record. Made by combining DNS name, address, and omg.lol top-level.",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *dnsRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dnsRecordResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan, and compute owner data

	var entry *omglol.DNSEntry
	if !plan.Priority.IsNull() {
		entry = omglol.NewDNSEntry(plan.Type.ValueString(), plan.Name.ValueString(), plan.Data.ValueString(), int(plan.TTL.ValueInt64()), int(plan.Priority.ValueInt64()))
	} else {
		entry = omglol.NewDNSEntry(plan.Type.ValueString(), plan.Name.ValueString(), plan.Data.ValueString(), int(plan.TTL.ValueInt64()))
	}

	// Set account settings
	record, err := r.client.CreateDNSRecord(plan.Address.ValueString(), *entry)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating settings",
			"Could not update settings, unexpected error: "+err.Error(),
		)
		return
	}

	plan.FQDN = types.StringValue(*record.Name + ".omg.lol")
	plan.CreatedAt = types.StringValue(*record.CreatedAt)
	plan.UpdatedAt = types.StringValue(*record.UpdatedAt)
	plan.ID = types.Int64Value(int64(*record.ID))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *dnsRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state dnsRecordResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := map[string]any{
		"ID": int(state.ID.ValueInt64()),
	}

	tflog.Warn(ctx, state.Address.ValueString())

	// Get refreshed DNS record from omg.lol
	record, err := r.client.FilterDNSRecord(state.Address.ValueString(), filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading DNS Record",
			"Could not read DNS record: "+err.Error(),
		)
		return
	}

	n := *record.Name
	lastIndex := strings.LastIndex(n, ".")
	address := n[lastIndex+1:]
	name := n[:lastIndex]

	// Overwrite record with refreshed state
	state.ID = types.Int64Value(int64(*record.ID))
	state.Type = types.StringValue(*record.Type)
	state.Address = types.StringValue(address)
	state.Name = types.StringValue(name)
	state.FQDN = types.StringValue(*record.Name + ".omg.lol")
	state.Data = types.StringValue(*record.Data)
	state.Priority = types.Int64Null()
	state.TTL = types.Int64Value(int64(*record.TTL))
	state.CreatedAt = types.StringValue(*record.CreatedAt)
	state.UpdatedAt = types.StringValue(*record.UpdatedAt)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dnsRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// // Retrieve values from plan
	// var plan dnsRecordResourceModel
	// diags := req.Plan.Get(ctx, &plan)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Generate API request body from plan, and compute owner data
	// settings := map[string]string{
	// 	"communication": plan.Communication.ValueString(),
	// 	"date_format":   plan.DateFormat.ValueString(),
	// }

	// // Set account settings
	// err := r.client.SetDNSRecord(settings)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error updating settings",
	// 		"Could not update settings, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	// plan.ID = types.StringValue("_")

	// // Set state to fully populated data
	// diags = resp.State.Set(ctx, plan)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dnsRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dnsRecordResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteDNSRecord(state.Address.ValueString(), int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting DNS Record",
			"Could not delete DNS record, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *dnsRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*omglol.Client)
}

func (r *dnsRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}