package provider

import (
	"context"
	"log"
	"os"

	"versa-networks.com/vclient"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// versaDirectorProviderModel maps provider schema data to a Go type.
type versaDirectorProviderModel struct {
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	Host              types.String `tfsdk:"host"`
	Port              types.String `tfsdk:"port"`
	OauthGrantType    types.String `tfsdk:"oauth_grant_type"`
	OauthClientID     types.String `tfsdk:"oauth_client_id"`
	OauthClientSecret types.String `tfsdk:"oauth_client_secret"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &versaDirectorProvider{}
)

// DataSources defines the data sources implemented in the provider.
func (p *versaDirectorProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAddressesDataSource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	log.Printf("New called .....\n")
	return func() provider.Provider {
		return &versaDirectorProvider{
			version: version,
		}
	}
}

// versaDIrectorProvider is the provider implementation.
type versaDirectorProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *versaDirectorProvider) Metadata(_ context.Context,
	_ provider.MetadataRequest, resp *provider.MetadataResponse) {

	resp.TypeName = "versadirector"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *versaDirectorProvider) Schema(_ context.Context,
	_ provider.SchemaRequest, resp *provider.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Interfact with versadirector",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Description: "Username for versadirector, May also be provided via VERSA_DIRECTOR_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for versadirector, May also be provided via VERSA_DIRECTOR_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"host": schema.StringAttribute{
				Description: "IP Address for versadirector, May also be provided via VERSA_DIRECTOR_HOST environment variable.",
				Optional:    true,
			},
			"port": schema.StringAttribute{
				Description: "Port for versadirector, May also be provided via VERSA_DIRECTOR_PORT environment variable.",
				Optional:    true,
			},
			"oauth_grant_type": schema.StringAttribute{
				Description: "Grant-Type for OAUTH2 authentication, May also be provided via VERSA_DIRECTOR_OAUTH_GRANT_TYPE environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"oauth_client_id": schema.StringAttribute{
				Description: "OAUTH2 Client-ID for authentication, May also be provided via VERSA_DIRECTOR_OAUTH_CLIENT_ID environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"oauth_client_secret": schema.StringAttribute{
				Description: "OAUTH2 Client-secret for authentication, May also be provided via VERSA_DIRECTOR_OAUTH_CLIENT_SECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *versaDirectorProvider) Configure(ctx context.Context,
	req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve provider data from configuration
	var config versaDirectorProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown versaDirector API Username",
			"The provider cannot create the versaDirector API client as there is an unknown configuration value for the versaDirector API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VERSA_DIRECTOR_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown versaDirector API Password",
			"The provider cannot create the versaDirector API client as there is an unknown configuration value for the versaDirector API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VERSA_DIRECTOR_PASSWORD environment variable.",
		)
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown versaDirector API Host",
			"The provider cannot create the versaDirector API client as there is an unknown configuration value for the versaDirector API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VERSA_DIRECTOR_HOST environment variable.",
		)
	}

	if config.Port.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			"Unknown versaDirector API Port",
			"The provider cannot create the versaDirector API client as there is an unknown configuration value for the versaDirector API port. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VERSA_DIRECTOR_PORT environment variable.",
		)
	}

	if config.OauthClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("oauth_client_id"),
			"Unknown versaDirector API OauthClientID",
			"The provider cannot create the versaDirector API client as there is an unknown configuration value for the versaDirector API oauth_client_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VERSA_DIRECTOR_OAUTH_CLIENT_ID environment variable.",
		)
	}

	if config.OauthClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("oauth_client_secret"),
			"Unknown versaDirector API OauthClientID",
			"The provider cannot create the versaDirector API client as there is an unknown configuration value for the versaDirector API oauth_client_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VERSA_DIRECTOR_CLIENT_OAUTH_SECRET environment variable.",
		)
	}

	if config.OauthGrantType.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("oauth_grant_type"),
			"Unknown versaDirector API OauthGrantType",
			"The provider cannot create the versaDirector API client as there is an unknown configuration value for the versaDirector API oauth_grant_type. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VERSA_DIRECTOR_OAUTH_GRANT_TYPE environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	username := os.Getenv("VERSA_DIRECTOR_USERNAME")
	password := os.Getenv("VERSA_DIRECTOR_PASSWORD")
	host := os.Getenv("VERSA_DIRECTOR_HOST")
	port := os.Getenv("VERSA_DIRECTOR_PORT")
	oauthGrantType := os.Getenv("VERSA_DIRECTOR_OAUTH_GRANT_TYPE")
	oauthClientId := os.Getenv("VERSA_DIRECTOR_OAUTH_CLIENT_ID")
	oauthClientSecret := os.Getenv("VERSA_DIRECTOR_OAUTH_CLIENT_SECRET")

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Port.IsNull() {
		port = config.Port.ValueString()
	}

	if !config.OauthGrantType.IsNull() {
		oauthGrantType = config.OauthGrantType.ValueString()
	}

	if !config.OauthClientID.IsNull() {
		oauthClientId = config.OauthClientID.ValueString()
	}

	if !config.OauthClientSecret.IsNull() {
		oauthClientSecret = config.OauthClientSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing HashiCups API Host",
			"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API host. "+
				"Set the host value in the configuration or use the VERSA_DIRECTOR_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing versaDirector API Username",
			"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API username. "+
				"Set the username value in the configuration or use the VERSA_DIRECTOR_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing versaDirector API Password",
			"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API password. "+
				"Set the password value in the configuration or use the VERSA_DIRECTOR_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if port == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			"Missing versaDirector API Port",
			"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API port. "+
				"Set the port value in the configuration or use the VERSA_DIRECTOR_PORT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if oauthGrantType == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("oauth_grant_type"),
			"Missing versaDirector API OauthGrantType",
			"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API oauthGrantType. "+
				"Set the oauth_grant_type value in the configuration or use the VERSA_DIRECTOR_OAUTH_GRANT_TYPE environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if oauthClientId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("oauth_client_id"),
			"Missing versaDirector API OauthClientId",
			"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API oauthClientId. "+
				"Set the oauth_client_id value in the configuration or use the VERSA_DIRECTOR_OAUTH_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if oauthClientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("oauth_client_secret"),
			"Missing versaDirector API OauthClientSecret",
			"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API oauthClientSecret. "+
				"Set the oauth_client_secret value in the configuration or use the VERSA_DIRECTOR_OAUTH_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new HashiCups client using the configuration values
	client, err := vclient.NewClient(&host, &username, &password, &port, &oauthClientId, &oauthClientSecret, &oauthGrantType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create versaDirector API Client",
			"An unexpected error occurred when creating the HashiCups API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"HashiCups Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines the resources implemented in the provider.
func (p *versaDirectorProvider) Resources(_ context.Context) []func() resource.Resource {
	log.Printf("Resources called .....\n")
	return []func() resource.Resource{
		NewAddressResource,
	}
}
