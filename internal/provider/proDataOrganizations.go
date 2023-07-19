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
	_ datasource.DataSource              = &organizationsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationsDataSource{}
)

// NewOrganizationsDataSource is a helper function to simplify the provider implementation.
func NewOrganizationsDataSource() datasource.DataSource {
	return &organizationsDataSource{}
}

// organizationsDataSource is the data source implementation.
type organizationsDataSource struct {
	client *vclient.Client
}

// Metadata returns the data source type name.
func (d *organizationsDataSource) Metadata(_ context.Context,
	req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_organizations"
}

// Schema defines the schema for the data source.
func (d *organizationsDataSource) Schema(_ context.Context,
	_ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organizations": schema.ListNestedAttribute{
				Description: "List of organizations in appliance.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Numeric identifer for the organization.",
							Computed:    true,
						},
						"uuid": schema.StringAttribute{
							Description: "Unique identifer for the organization.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the organization.",
							Computed:    true,
						},
						"subscription_plan": schema.StringAttribute{
							Description: "Subscription plan for organization.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *organizationsDataSource) Configure(_ context.Context,
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
func (d *organizationsDataSource) Read(ctx context.Context,
	req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Info(ctx, "Read organizations data")

	var state organizationsDataSourceList

	organizations, err := d.client.GetAllOrganizations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Organizations,", err.Error())
		return
	}
	for _, val := range organizations {
		curOrganization := organizationsData{
			Id:               types.Int64Value(int64(val.Id)),
			Uuid:             types.StringValue(val.UUID),
			Name:             types.StringValue(val.Name),
			SubscriptionPlan: types.StringValue(val.SubscriptionPlan),
		}
		tflog.Info(ctx, "Updating org data")
		state.Organizations = append(state.Organizations, curOrganization)

	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// organizationsDataSourceList maps the data source schema data.
type organizationsDataSourceList struct {
	Organizations []organizationsData `tfsdk:"organizations"`
}

// organizationsData maps Organizations schema data.
type organizationsData struct {
	Id               types.Int64  `tfsdk:"id"`
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	SubscriptionPlan types.String `tfsdk:"subscription_plan"`
}
