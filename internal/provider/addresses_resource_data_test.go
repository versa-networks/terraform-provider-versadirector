package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "versadirector_addresses" "test" {
  device_name = "Branch-1"
  organization_name = "Customer-1"
  address = [
    {
    	name = "versa-networks-addresses-1"
      	fqdn = "versa-networks.com/addresses-1"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first address item
					//resource.TestCheckResourceAttr("versadirector_addresses.test", "address.0.name", "versa-networks-addresses"),
					//resource.TestCheckResourceAttr("versadirector_addresses.test", "address.1.name", "juniper-networks-addresses"),
					//resource.TestCheckResourceAttr("versadirector_addresses.test", "address.2.name", "cisco-networks-addresses"),
					resource.TestCheckResourceAttrSet("versadirector_addresses.test", "address.0.name"),
					//resource.TestCheckResourceAttrSet("versadirector_addresses.test", "address.1.name"),
					//resource.TestCheckResourceAttrSet("versadirector_addresses.test", "address.2.name"),
					resource.TestCheckResourceAttrSet("versadirector_addresses.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "versadirector_addresses.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "versadirector_addresses" "test" {
  device_name = "Branch-1"
  organization_name = "Customer-1"
  address = [
    {
    	name = "versa-networks-addresses-1"
      	fqdn = "versa-networks.com/addresses-2"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first coffee item has Computed attributes updated.
					//resource.TestCheckResourceAttr("versadirector_addresses.test", "address.0.name", "versa-networks-addresses"),
					//resource.TestCheckResourceAttr("versadirector_addresses.test", "address.1.name", "cisco-networks-addresses"),
					//resource.TestCheckResourceAttr("versadirector_addresses.test", "address.2.name", "juniper-networks-addresses"),
					resource.TestCheckResourceAttrSet("versadirector_addresses.test", "address.0.name"),
					//resource.TestCheckResourceAttrSet("versadirector_addresses.test", "address.1.name"),
					//resource.TestCheckResourceAttrSet("versadirector_addresses.test", "address.2.name"),
					resource.TestCheckResourceAttrSet("versadirector_addresses.test", "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
