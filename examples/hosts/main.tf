terraform {
  required_providers {
    onos = {
      source = "ctjnkns/onos"
    }
  }
}

provider "onos" {
  host     = "http://localhost:8181/onos/v1"
  username = "onos"
  password = "rocks"
}

data "onos_hosts" "mininet" {}

output "mininet_hosts" {
  value = data.onos_hosts.mininet
}