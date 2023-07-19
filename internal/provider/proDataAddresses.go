package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"versa-networks.com/vclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &addressesDataSource{}
	_ datasource.DataSourceWithConfigure = &addressesDataSource{}
)

// NewAddressesDataSource is a helper function to simplify the provider implementation.
func NewAddressesDataSource() datasource.DataSource {
	return &addressesDataSource{}
}

// addressesDataSource is the data source implementation.
type addressesDataSource struct {
	client *vclient.Client
}

// Metadata returns the data source type name.
func (d *addressesDataSource) Metadata(_ context.Context,
	req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_addresses"
}

// Schema defines the schema for the data source.
func (d *addressesDataSource) Schema(ctx context.Context,
	_ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier name to be configured.",
				Computed:    true,
			},
			"device_name": schema.StringAttribute{
				Description: "Device name to be configured.",
				Required:    true,
			},
			"organization_name": schema.StringAttribute{
				Description: "Organization name for the device to be configured.",
				Required:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the configuration.",
				Computed:    true,
			},
			"address": schema.ListNestedAttribute{
				Description: "List of configured addresses.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the address object.",
							Optional:    true,
						},
						"fqdn": schema.StringAttribute{
							Description: "FQDN for the address object.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *addressesDataSource) Configure(_ context.Context,
	req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*vclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *vclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *addressesDataSource) Read(ctx context.Context,
	req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config addressesDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceName := config.DeviceName.ValueString()
	if len(deviceName) <= 0 {
		deviceName = os.Getenv("VERSA_VOS_DEVICE_NAME")
		if len(deviceName) <= 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("deviceName"),
				"Missing deviceName for versaDirector API",
				"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API deviceName. "+
					"Set the username value in the configuration or use the VERSA_VOS_DEVICE_NAME environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}
	}

	organizationName := config.OrganizationName.ValueString()
	if len(organizationName) <= 0 {
		organizationName = os.Getenv("VERSA_VOS_ORGANIZATION_NAME")
		if len(deviceName) <= 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("deviceName"),
				"Missing deviceName for versaDirector API",
				"The provider cannot create the versaDirector API client as there is a missing or empty value for the versaDirector API deviceName. "+
					"Set the username value in the configuration or use the VERSA_VOS_ORGANIZATION_NAME environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}
	}

	tflog.Debug(ctx, "DATA-READ: Get Addresses for Device: "+deviceName+
		" Organization: "+organizationName)

	addrData, err := d.client.GetDeviceOrganizationAddresses(ctx, deviceName, organizationName)
	if err != nil {
		tflog.Error(ctx, "Failed to get addresses for device "+deviceName+" Organization "+organizationName)
	} else {
		if len(deviceName) > 0 {
			config.DeviceName = types.StringValue(deviceName)
		}
		if len(organizationName) > 0 {
			config.OrganizationName = types.StringValue(organizationName)
		}
		for _, val := range addrData.Addresses {
			configAddrData := addressDataModel{
				Name: types.StringValue(val.Name),
				FQDN: types.StringValue(val.FQDN),
			}
			config.Address = append(config.Address, configAddrData)
		}
	}
	config.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	config.ID = types.StringValue("1")

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}

// orderResourceModel maps the resource schema data.
type addressesDataModel struct {
	ID               types.String       `tfsdk:"id"`
	DeviceName       types.String       `tfsdk:"device_name"`
	OrganizationName types.String       `tfsdk:"organization_name"`
	LastUpdated      types.String       `tfsdk:"last_updated"`
	Address          []addressDataModel `tfsdk:"address"`
}

// addressItem
type addressDataModel struct {
	Name types.String `tfsdk:"name"`
	FQDN types.String `tfsdk:"fqdn"`
}
