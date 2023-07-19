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

data "versadirector_addresses" "vos" {
  device_name       = "DEVICE_NAME"
  organization_name = "ORAGANIZATION"
}

output "vos_addresses" {
  value = data.versadirector_addresses.vos
}
