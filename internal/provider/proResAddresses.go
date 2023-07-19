package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"versa-networks.com/vclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &addressResource{}
	_ resource.ResourceWithConfigure   = &addressResource{}
	_ resource.ResourceWithImportState = &addressResource{}
)

// NewAddressesResource is a helper function to simplify the provider implementation.
func NewAddressResource() resource.Resource {
	return &addressResource{}
}

// addressResource is the resource implementation.
type addressResource struct {
	client *vclient.Client
}

func (r *addressResource) ImportState(ctx context.Context,
	req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata returns the resource type name.
func (r *addressResource) Metadata(_ context.Context,
	req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_addresses"
}

// Schema defines the schema for the resource.
func (r *addressResource) Schema(_ context.Context,
	_ resource.SchemaRequest, resp *resource.SchemaResponse) {

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
				Description: "List of addresses to be configured.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the address object.",
							Required:    true,
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

// Configure adds the provider configured client to the resource.
func (r *addressResource) Configure(_ context.Context,
	req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *addressResource) Create(ctx context.Context,
	req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan addressesResourceModel

	tflog.Debug(ctx, "CREATE Address request received")

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addrListData := vclient.DevOjectsAddressListData{}
	addrList := &addrListData.AddrList
	addrList.Count = len(plan.Address)
	addrList.DeviceName = plan.DeviceName.ValueString()
	addrList.OrganizationName = plan.OrganizationName.ValueString()

	for _, val := range plan.Address {
		address := vclient.DevObjectAddress{
			Name: val.Name.ValueString(),
			FQDN: val.FQDN.ValueString(),
		}
		addrList.Addresses = append(addrList.Addresses, address)
	}

	r.client.CreateDevOrgServiceObjAddresses(ctx, addrListData)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.ID = types.StringValue("1")

	tflog.Debug(ctx, "RESOURCE Create for Device: "+addrList.DeviceName+
		" Organization: "+addrList.OrganizationName)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "CREATE Address request completed")
}

// Read refreshes the Terraform state with the latest data.
func (r *addressResource) Read(ctx context.Context,
	req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Debug(ctx, "READ Address request received")

	// Get current state
	var state addressesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceName := state.DeviceName.ValueString()
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

	organizationName := state.OrganizationName.ValueString()
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

	//fmt.Printf("RES-READ device %v Org %v\n", deviceName, organizationName)
	addrData, err := r.client.GetDeviceOrganizationAddresses(ctx, deviceName, organizationName)
	if err != nil {
		tflog.Error(ctx, "Failed to get addresses for device "+deviceName+" Organization "+organizationName)
	} else {
		state.DeviceName = types.StringValue(deviceName)
		state.OrganizationName = types.StringValue(organizationName)
		for _, val := range addrData.Addresses {
			configAddrData := addressDataModel{
				Name: types.StringValue(val.Name),
				FQDN: types.StringValue(val.FQDN),
			}
			//fmt.Printf("Adding address %v\n", configAddrData)
			//state.Address = append(state.Address, addressItemModel(configAddrData))
			added := false
			for stKey, stVal := range state.Address {
				if stVal.Name == configAddrData.Name {
					added = true
					state.Address[stKey] = addressItemModel(configAddrData)
					break
				}
			}
			if added == false {
				state.Address = append(state.Address, addressItemModel(configAddrData))
			}
		}
	}
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	state.ID = types.StringValue("1")

	tflog.Debug(ctx, "RESOURCE Read for Device: "+deviceName+
		" Organization: "+organizationName)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "READ Address request completed")
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *addressResource) Update(ctx context.Context,
	req resource.UpdateRequest, resp *resource.UpdateResponse) {

	tflog.Debug(ctx, "UPDATE Address request received")

	// Retrieve values from plan
	var plan addressesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addrListData := vclient.DevOjectsAddressListData{}
	addrList := &addrListData.AddrList
	addrList.Count = len(plan.Address)
	addrList.DeviceName = plan.DeviceName.ValueString()
	addrList.OrganizationName = plan.OrganizationName.ValueString()

	for _, val := range plan.Address {
		address := vclient.DevObjectAddress{
			Name: val.Name.ValueString(),
			FQDN: val.FQDN.ValueString(),
		}
		addrList.Addresses = append(addrList.Addresses, address)
	}

	r.client.UpdateDevOrgServiceObjAddresses(ctx, addrListData)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.ID = types.StringValue("1")

	tflog.Debug(ctx, "RESOURCE Update for Device: "+addrList.DeviceName+
		" Organization: "+addrList.OrganizationName)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *addressResource) Delete(ctx context.Context,
	req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Debug(ctx, "DELETE Address request received")

	var state addressesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addrListData := vclient.DevOjectsAddressListData{}
	addrList := &addrListData.AddrList
	addrList.Count = len(state.Address)
	addrList.DeviceName = state.DeviceName.ValueString()
	addrList.OrganizationName = state.OrganizationName.ValueString()

	for _, val := range state.Address {
		address := vclient.DevObjectAddress{
			Name: val.Name.ValueString(),
			FQDN: val.FQDN.ValueString(),
		}
		addrList.Addresses = append(addrList.Addresses, address)
	}
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	state.ID = types.StringValue("1")

	r.client.DeleteDevOrgServiceObjAddresses(ctx, addrListData)

	tflog.Debug(ctx, "RESOURCE Delete for Device: "+addrList.DeviceName+
		" Organization: "+addrList.OrganizationName)
}

// orderResourceModel maps the resource schema data.
type addressesResourceModel struct {
	ID               types.String       `tfsdk:"id"`
	DeviceName       types.String       `tfsdk:"device_name"`
	OrganizationName types.String       `tfsdk:"organization_name"`
	LastUpdated      types.String       `tfsdk:"last_updated"`
	Address          []addressItemModel `tfsdk:"address"`
}

// addressItem
type addressItemModel struct {
	Name types.String `tfsdk:"name"`
	FQDN types.String `tfsdk:"fqdn"`
}
