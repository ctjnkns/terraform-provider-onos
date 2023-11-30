package provider

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &onosProvider{}
)

// onosProviderModel maps provider schema data to a Go type.
type onosProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type onosClient struct {
	HTTPClient *http.Client
	Host       string
	Username   string
	Password   string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &onosProvider{
			version: version,
		}
	}
}

// onosProvider is the provider implementation.
type onosProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *onosProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "onos"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *onosProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *onosProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config onosProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown onos API Host",
			"The provider cannot create the onos API client as there is an unknown configuration value for the onos API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the onos_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown onos API Username",
			"The provider cannot create the onos API client as there is an unknown configuration value for the onos API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the onos_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown onos API Password",
			"The provider cannot create the onos API client as there is an unknown configuration value for the onos API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the onos_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("ONOS_HOST")
	username := os.Getenv("ONOS_USERNAME")
	password := os.Getenv("ONOS_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing onos API Host",
			"The provider cannot create the onos API client as there is a missing or empty value for the onos API host. "+
				"Set the host value in the configuration or use the ONOS_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing onos API Username",
			"The provider cannot create the onos API client as there is a missing or empty value for the onos API username. "+
				"Set the username value in the configuration or use the ONOS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing onos API Password",
			"The provider cannot create the onos API client as there is a missing or empty value for the onos API password. "+
				"Set the password value in the configuration or use the ONOS_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Onos client using the configuration values
	//client := &http.Client{Timeout: 10 * time.Second}

	// Make the Onos client available during DataSource and Resource
	// type Configure methods.

	client := &http.Client{Timeout: 10 * time.Second}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *onosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFlowsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *onosProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewIntentResource,
	}
}
