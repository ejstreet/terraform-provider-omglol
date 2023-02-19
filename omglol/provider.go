package omglol

import (
	"context"
	"os"

	"github.com/ejstreet/omglol-client-go/omglol"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &omglolProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &omglolProvider{}
}

// omglolProvider is the provider implementation.
type omglolProvider struct{}

// Metadata returns the provider type name.
func (p *omglolProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "omglol"
}

// Schema defines the provider-level schema for configuration data.
func (p *omglolProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The omg.lol provider provides utilities for working with resources in omg.lol.",
		Attributes: map[string]schema.Attribute{
			"api_host": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "This variable is not required, and only useful for development purposes. Default value is `https://api.omg.lol`. Pass this variable in the provider configuration, or alternatively set the `OMGLOL_API_HOST` environment variable.",
			},
			"user_email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Pass this variable in the provider configuration, or alternatively set the `OMGLOL_USER_EMAIL` environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Pass this variable in the provider configuration, or alternatively set the `OMGLOL_API_KEY` environment variable. As this is a sensitive variable, it is recommended to set it as an environment variable.",
			},
		},
	}
}

// omglolProviderModel maps provider schema data to a Go type.
type omglolProviderModel struct {
	APIHost   types.String `tfsdk:"api_host"`
	APIKey    types.String `tfsdk:"api_key"`
	UserEmail types.String `tfsdk:"user_email"`
}

// Configure prepares a omglol API client for data sources and resources.
func (p *omglolProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config omglolProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.APIHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_host"),
			"Unknown omglol API Host",
			"The provider cannot create the omg.lol API client as there is an unknown configuration value for the omg.lol API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMGLOL_HOST environment variable.",
		)
	}

	if config.UserEmail.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("user_email"),
			"Unknown omglol API User_email",
			"The provider cannot create the omg.lol API client as there is an unknown configuration value for the omg.lol API user_email. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMGLOL_USER_EMAIL environment variable.",
		)
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown omglol API Api_key",
			"The provider cannot create the omg.lol API client as there is an unknown configuration value for the omg.lol API api_key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMGLOL_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("OMGLOL_API_HOST")
	api_key := os.Getenv("OMGLOL_API_KEY")
	user_email := os.Getenv("OMGLOL_USER_EMAIL")

	if !config.APIHost.IsNull() {
		host = config.APIHost.ValueString()
	} else {
		host = "https://api.omg.lol"
	}

	if !config.APIKey.IsNull() {
		api_key = config.APIKey.ValueString()
	}

	if !config.UserEmail.IsNull() {
		user_email = config.UserEmail.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if user_email == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("user_email"),
			"Missing omglol API User Email address",
			"The provider cannot create the omglol API client as there is a missing or empty value for the omglol API user email address. "+
				"Set the user_email value in the configuration or use the OMGLOL_USER_EMAIL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing omglol API key",
			"The provider cannot create the omglol API client as there is a missing or empty value for the omglol API key. "+
				"Set the api_key value in the configuration or use the OMGLOL_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "host", host)
	ctx = tflog.SetField(ctx, "user_email", user_email)
	ctx = tflog.SetField(ctx, "api_key", api_key)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_key")

	tflog.Debug(ctx, "Creating omg.lol client")

	// Create a new omglol client using the configuration values
	client, err := omglol.NewClient(user_email, api_key, host)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create omglol API Client",
			"An unexpected error occurred when creating the omglol API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"omglol Client Error: "+err.Error(),
		)
		return
	}

	// Make the omglol client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured omg.lol client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *omglolProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountInfoDataSource,
		NewDnsRecordsDataSource,
		NewPURLsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *omglolProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccountSettingsResource,
		NewDNSRecordResource,
		NewPURLResource,
	}
}
