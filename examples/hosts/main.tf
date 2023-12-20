terraform {
  required_providers {
    onos = {
      source = "hashicorp.com/edu/onos"
    }
  }
}

provider "onos" {
  host     = "http://localhost:8181/onos/v1"
  username = "onos"
  password = "rocks"
}

data "onos_hosts" "edu" {}

output "edu_hosts" {
  value = data.onos_hosts.edu
}