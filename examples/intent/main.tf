terraform {
  required_providers {
    onos = {
      source = "hashicorp.com/edu/onos"
    }
  }
}

provider "onos" {
  host     = "http://localhost:8181"
  username = "onos"
  password = "rocks"
}

resource "onos_intent" "edu" {
  items = [{
    coffee = {
      id = 3
    }
    quantity = 2
    }, {
    coffee = {
      id = 1
    }
    quantity = 2
    }
  ]
}

output "edu_intent" {
  value = onos_intent.edu
}
