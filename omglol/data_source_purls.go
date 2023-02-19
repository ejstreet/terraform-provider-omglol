package omglol

import (
	"context"

	"github.com/ejstreet/omglol-client-go/omglol"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &pURLsDataSource{}
	_ datasource.DataSourceWithConfigure = &pURLsDataSource{}
)

func NewPURLsDataSource() datasource.DataSource {
	return &pURLsDataSource{}
}

type pURLsDataSource struct {
	client *omglol.Client
}

// Configure adds the provider configured client to the data source.
func (d *pURLsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*omglol.Client)
}

func (d *pURLsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_purls"
}

func (d *pURLsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List all PURLs for a given omg.lol address.",
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The omg.lol address to read the purls from.",
			},
			"purls": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A list of all the PURLs for the given address.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the PURL. The name field is how you will access your designated URL.",
						},
						"url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The url that is pointed to.",
						},
						"listed": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Returns `true` if listed on your `address`.url.lol page.",
						},
						"counter": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of time a PURL has been accessed.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique ID of the PURL. Can be used for imports.",
						},
					},
				},
			},
		},
	}
}

type pURLDataSourceModel struct {
	Name    types.String `tfsdk:"name"`
	URL     types.String `tfsdk:"url"`
	Listed  types.Bool   `tfsdk:"listed"`
	Counter types.Int64  `tfsdk:"counter"`
	ID      types.String `tfsdk:"id"`
}

type pURLsDataSourceModel struct {
	Address types.String          `tfsdk:"address"`
	PURLs   []pURLDataSourceModel `tfsdk:"purls"`
}

func (d *pURLsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state pURLsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pURLs, err := d.client.ListPersistentURLs(state.Address.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PURLs",
			err.Error(),
		)
		return
	}

	for _, purl := range *pURLs {

		var counter int64
		if purl.Counter != nil {
			counter = *purl.Counter
		} else {
			counter = 0
		}

		p := pURLDataSourceModel{
			ID:      types.StringValue(state.Address.ValueString() + "_" + purl.Name),
			Name:    types.StringValue(purl.Name),
			URL:     types.StringValue(purl.URL),
			Counter: types.Int64Value(counter),
		}

		state.PURLs = append(state.PURLs, p)

	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
