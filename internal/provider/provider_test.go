package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the versa client is properly configured.
	providerConfig = `
provider "versadirector" {
	username            = "Administrator"
	password            = "Versa123#"
	host                = "10.40.73.242"
	port                = "9182"
	oauth_grant_type    = "password"
	oauth_client_id     = "CA736092A7221051EA93B4447A259744"
	oauth_client_secret = "6bafb4e78e909775a377eedf022610a6"
}
data "versadirector_addresses" "test" {
	device_name       = "Branch-1"
	organization_name = "ACME"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"versadirector": providerserver.NewProtocol6WithError(New("test")()),
	}
)
