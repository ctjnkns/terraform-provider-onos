terraform {
  required_providers {
    onos = {
      source = "hashicorp.com/ctjnkns/onos"
    }
  }
}

provider "onos" {}

data "onos_flows" "example" {}
