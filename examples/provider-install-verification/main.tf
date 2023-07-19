terraform {
  required_providers {
    versadirector = {
      source = "versa-networks.com/versa-networks/versadirector"
    }
  }
}

provider "versadirector" {}

data "versadirector_config" "example" {}

