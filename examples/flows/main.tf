terraform {
  required_providers {
    onos = {
      source = "hashicorp.com/edu/onos"
    }
  }
}

provider "onos" {
  host     = "http://localhost:8181"
  username = "education"
  password = "test123"
}

data "onos_flows" "edu" {}