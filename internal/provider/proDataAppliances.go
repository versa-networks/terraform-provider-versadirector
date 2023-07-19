package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"versa-networks.com/vclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &appliancesDataSource{}
	_ datasource.DataSourceWithConfigure = &appliancesDataSource{}
)

// NewAppliancesDataSource is a helper function to simplify the provider implementation.
func NewAppliancesDataSource() datasource.DataSource {
	return &appliancesDataSource{}
}

// appliancesDataSource is the data source implementation.
type appliancesDataSource struct {
	client *vclient.Client
}

// Metadata returns the data source type name.
func (d *appliancesDataSource) Metadata(_ context.Context,
	req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_appliances"
}

// Schema defines the schema for the data source.
func (d *appliancesDataSource) Schema(ctx context.Context,
	_ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	tflog.Info(ctx, "Schema called....")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"appliance_name": schema.StringAttribute{
				Description: "Identifier for this address item.",
				Required:    true,
			},
			"organization_name": schema.StringAttribute{
				Description: "Org Identifier for this address item.",
				Required:    true,
			},
			"appliances": schema.ListNestedAttribute{
				Description: "List of appliances.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "Unique identifier for the appliance.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the appliance.",
							Computed:    true,
						},
						"services_status": schema.StringAttribute{
							Description: "Health status of services in appliance.",
							Computed:    true,
						},
						"overall_status": schema.StringAttribute{
							Description: "Overall health status of the appliance.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *appliancesDataSource) Configure(_ context.Context,
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
func (d *appliancesDataSource) Read(ctx context.Context,
	req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state appliancesDataSourceList

	appliances, err := d.client.GetAllAppliances(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading FlexVNF Appliances,", err.Error())
		return
	}

	for i := 0; i < appliances.TotalCount; i++ {
		curAppliance := applianceData{
			Uuid:           types.StringValue(appliances.Appliances[i].UUID),
			Name:           types.StringValue(appliances.Appliances[i].Name),
			ServicesStatus: types.StringValue(appliances.Appliances[i].ServicesStatus),
			OverallStatus:  types.StringValue(appliances.Appliances[i].OverallStatus),
		}
		state.Appliances = append(state.Appliances, curAppliance)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// appliancesDataSourceList maps the data source schema data.
type appliancesDataSourceList struct {
	ApplianceName   types.String    `tfsdk:"appliance_name"`
	OrganizatioName types.String    `tfsdk:"organization_name"`
	Appliances      []applianceData `tfsdk:"appliances"`
}

// applianceData maps appliances schema data.
type applianceData struct {
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	ServicesStatus types.String `tfsdk:"services_status"`
	OverallStatus  types.String `tfsdk:"overall_status"`
}
