terraform {
  required_providers {
    versadirector = {
      source = "versa-networks.com/versa-networks/versadirector"
    }
  }
}


provider "versadirector" {
  username            = "username"
  password            = "password"
  host                = "10.20.30.40"
  port                = "9182"
  #always to be set oauth_grant_type = "password"
  oauth_grant_type    = "password"
  oauth_client_id     = "XXXXXXXXXXXXXXXX"
  oauth_client_secret = "YYYYYYYYYYYYYYYY"
}

resource "versadirector_addresses" "vos_addresses" {
  device_name       = "devicename"
  organization_name = "orgname"
  address = [
    {
      name = "versa-networks-addresses1"
      fqdn = "versa-networks.com/addresses"
    },
   {
      name = "versa-networks-addresses2"
      fqdn = "versa-networks.com/addresses"
    },
    {
      name = "versa-networks-addresses3"
      fqdn = "versa-networks.com/addresses"
    }
  ]
}

output "versa_addresses_info" {
  value = versadirector_addresses.vos_addresses
}
