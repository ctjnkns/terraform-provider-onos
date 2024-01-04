#Manage example intent
resource "onos_intent" "example" {
  intent = {
    appid    = "org.onosproject.cli"
    key      = "0x100006"
    type     = "HostToHostIntent"
    priority = 100
    one      = "00:00:00:00:00:01/None"
    two      = "00:00:00:00:00:02/None"
  }
}