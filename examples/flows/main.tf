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

data "onos_flows" "edu" {}

output "edu_flows" {
  value = data.onos_flows.edu
}