terraform {
  required_providers {
    onos = {
      source = "ctjnkns/onos"
    }
  }
}

provider "onos" {}

data "onos_flows" "example" {}
