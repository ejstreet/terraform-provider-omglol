package omglol

import (
	"context"
	"strings"

	"github.com/ejstreet/omglol-client-go/omglol"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dnsRecordsDataSource{}
	_ datasource.DataSourceWithConfigure = &dnsRecordsDataSource{}
)

func NewDnsRecordsDataSource() datasource.DataSource {
	return &dnsRecordsDataSource{}
}

type dnsRecordsDataSource struct {
	client *omglol.Client
}

// Configure adds the provider configured client to the data source.
func (d *dnsRecordsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*omglol.Client)
}

func (d *dnsRecordsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_records"
}

func (d *dnsRecordsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all DNS records for a given omg.lol address.",
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The omg.lol address to read the records from.",
			},
			"records": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A list of all the DNS records for the given address.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The record type.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The prefix attached before the address. `@` represents the top level.",
						},
						"data": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The data entered into the record.",
						},
						"ttl": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The Time-To-Live (TTL) of the record.",
						},
						"priority": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The priority of the record. Only applies to MX records.",
						},
						"fqdn": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The fully qualified domain name of the record. Made by combining DNS name, address, and omg.lol top-level.",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type dnsRecordDataSourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Name      types.String `tfsdk:"name"`
	Data      types.String `tfsdk:"data"`
	Priority  types.Int64  `tfsdk:"priority"`
	TTL       types.Int64  `tfsdk:"ttl"`
	FQDN      types.String `tfsdk:"fqdn"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

type dnsRecordsDataSourceModel struct {
	Address types.String               `tfsdk:"address"`
	Records []dnsRecordDataSourceModel `tfsdk:"records"`
}

func (d *dnsRecordsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dnsRecordsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsRecords, err := d.client.ListDNSRecords(state.Address.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read DNS Records",
			err.Error(),
		)
		return
	}

	for _, record := range *dnsRecords {

		r := dnsRecordDataSourceModel{
			ID:        types.Int64Value(int64(*record.ID)),
			Type:      types.StringValue(*record.Type),
			Data:      types.StringValue(*record.Data),
			TTL:       types.Int64Value(int64(*record.TTL)),
			FQDN:      types.StringValue((*record.Name + ".omg.lol")),
			CreatedAt: types.StringValue(*record.CreatedAt),
			UpdatedAt: types.StringValue(*record.UpdatedAt),
		}

		if strings.Contains(*record.Name, ".") {
			r.Name = types.StringValue(strings.Split(*record.Name, ".")[0])
		} else {
			r.Name = types.StringValue("@")
		}

		if *record.Type == "MX" {
			r.Priority = types.Int64Value(int64(*record.Priority))
		} else {
			r.Priority = types.Int64Null()
		}

		state.Records = append(state.Records, r)

	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
