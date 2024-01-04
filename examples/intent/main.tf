terraform {
  required_providers {
    onos = {
      source = "ctjnkns/onos"
    }
  }
  required_version = ">= 1.1.0"
}

provider "onos" {
  host     = "http://localhost:8181/onos/v1"
  username = "onos"
  password = "rocks"
}


resource "onos_intent" "h1-to-h2" {
  intent = {
    appid    = "org.onosproject.cli"
    key      = "0x100006"
    type     = "HostToHostIntent"
    priority = 100
    one      = "00:00:00:00:00:01/None"
    two      = "00:00:00:00:00:02/None"
  }
}

output "h1-to-h2_intent" {
  value = onos_intent.h1-to-h2
}

/*
resource "onos_intent" "h2-to-h3" {
  intent = {
    appid    = "org.onosproject.cli"
    key      = "0x100007"
    type     = "HostToHostIntent"
    priority = 100
    one      = "00:00:00:00:00:02/None"
    two      = "00:00:00:00:00:03/None"
  }
}

output "h2-to-h3_intent" {
  value = onos_intent.h2-to-h3
}
*/

