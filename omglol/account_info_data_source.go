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
	_ datasource.DataSource              = &accountInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &accountInfoDataSource{}
)

func NewAccountInfoDataSource() datasource.DataSource {
	return &accountInfoDataSource{}
}

type accountInfoDataSource struct {
	client *omglol.Client
}

// Configure adds the provider configured client to the data source.
func (d *accountInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*omglol.Client)
}

func (d *accountInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_info"
}

func (d *accountInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve account info.",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The email address associated with the account.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name associated with the account.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The RFC 3339 representation of the time that the account was created. This can be used in conjunction with the [formatdate](https://developer.hashicorp.com/terraform/language/functions/formatdate) function.",
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type accountInfoDataSourceModel struct {
	Email   types.String `tfsdk:"email"`
	Name    types.String `tfsdk:"name"`
	Created types.String `tfsdk:"created"`
	ID      types.String `tfsdk:"id"`
}

func (d *accountInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	accountInfo, err := d.client.GetAccountInfo()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Account Info",
			err.Error(),
		)
		return
	}

	state := accountInfoDataSourceModel{
		Email:   types.StringValue(accountInfo.Email),
		Name:    types.StringValue(accountInfo.Name),
		Created: types.StringValue(accountInfo.Created.Iso8601Time),
		ID:      types.StringValue("_"),
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
