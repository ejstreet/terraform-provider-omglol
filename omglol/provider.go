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
		Attributes: map[string]schema.Attribute{
			"api_endpoint": schema.StringAttribute{
				Optional: true,
			},
			"user_email": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// omglolProviderModel maps provider schema data to a Go type.
type omglolProviderModel struct {
	Api_endpoint types.String `tfsdk:"api_endpoint"`
	User_email   types.String `tfsdk:"user_email"`
	Api_key      types.String `tfsdk:"api_key"`
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

	if config.Api_endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_endpoint"),
			"Unknown omglol API Api_endpoint",
			"The provider cannot create the omg.lol API client as there is an unknown configuration value for the omg.lol API api_endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMGLOL_API_ENDPOINT environment variable.",
		)
	}

	if config.User_email.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("user_email"),
			"Unknown omglol API User_email",
			"The provider cannot create the omg.lol API client as there is an unknown configuration value for the omg.lol API user_email. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OMGLOL_USER_EMAIL environment variable.",
		)
	}

	if config.Api_key.IsUnknown() {
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

	api_endpoint := os.Getenv("OMGLOL_API_ENDPOINT")
	api_key := os.Getenv("OMGLOL_API_KEY")
	user_email := os.Getenv("OMGLOL_USER_EMAIL")

	if !config.Api_endpoint.IsNull() {
		api_endpoint = config.Api_endpoint.ValueString()
	} else {
		api_endpoint = "https://api.omg.lol"
	}

	if !config.Api_key.IsNull() {
		api_key = config.Api_key.ValueString()
	}

	if !config.User_email.IsNull() {
		user_email = config.User_email.ValueString()
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

	// Create a new omglol client using the configuration values
	client, err := omglol.NewClient(user_email, api_key, api_endpoint)
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
}

// DataSources defines the data sources implemented in the provider.
func (p *omglolProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountInfoDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *omglolProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
