package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAddressesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				//Config: providerConfig + `data "versadirector_addresses" "test" {}`,
				Config: providerConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.versadirector_addresses.test", "device_name", "Branch-1"),
					resource.TestCheckResourceAttr("data.versadirector_addresses.test", "organization_name", "ACME"),
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.versadirector_addresses.test", "address.0.name", "cisco-systems-addresses"),
					resource.TestCheckResourceAttr("data.versadirector_addresses.test", "address.1.name", "jnupier-networks-addresses"),
				),
			},
		},
	})
}
