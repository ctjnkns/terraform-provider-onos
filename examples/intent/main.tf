terraform {
  required_providers {
    onos = {
      source = "hashicorp.com/edu/onos"
    }
  }
  required_version = ">= 1.1.0"
}

provider "onos" {
  host     = "http://localhost:8181/onos/v1"
  username = "onos"
  password = "rocks"
}

resource "onos_intent" "edu" {
  intent = {
    appid    = "org.onosproject.cli"
    key      = "0x100006"
    type     = "HostToHostIntent"
    priority = 100
    one      = "00:00:00:00:00:02/None"
    two      = "00:00:00:00:00:88/None"
  }
}


resource "onos_intent" "io" {
  intent = {
    appid    = "org.onosproject.cli"
    key      = "0x100005"
    type     = "HostToHostIntent"
    priority = 100
    one      = "00:00:00:00:00:01/None"
    two      = "00:00:00:00:00:02/None"
  }
}

output "edu_intent" {
  value = onos_intent.edu
}

output "io_intent" {
  value = onos_intent.io
}