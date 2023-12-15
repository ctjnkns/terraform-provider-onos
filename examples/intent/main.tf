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
    intents = [
        {
            "type": "HostToHostIntent",
            "id": "0x100005",
            "key": "0x100005",
            "appId": "org.onosproject.cli",
            "resources": [
                "00:00:00:00:00:01/None",
                "00:00:00:00:00:03/None"
            ],
            "state": "FAILED"
        },
        {
            "type": "HostToHostIntent",
            "id": "0x100000",
            "key": "0x100000",
            "appId": "org.onosproject.cli",
            "resources": [
                "00:00:00:00:00:01/None",
                "00:00:00:00:00:02/None"
            ],
            "state": "INSTALLED"
        }
    ]
}

output "edu_intent" {
  value = onos_intent.edu
}
