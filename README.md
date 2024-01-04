# ONOS Terraform Provider

This terraform provider is based on the templates provied in [Terraform Provider Scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding-framework) and the [Implement a provider with the Terraform Plugin Framework](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider) documentation.

A basic understanding of Networking, [Mininet](https://github.com/mininet/mininet/wiki/Documentation), and [ONOS](https://opennetworking.org/onos/) is assumed.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

### Basic Usage
#### Intents
The intents resource requires the fields: appid, key, type, priority, one, and two. The id field is computed by onos and added to the state after receiving the value. The last_updated field is computed during terraform executions.

Configuration:
```hcl
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
```

Output:
```shell
$ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # onos_intent.h1-to-h2 will be created
  + resource "onos_intent" "h1-to-h2" {
      + id           = (known after apply)
      + intent       = {
          + appid    = "org.onosproject.cli"
          + id       = (known after apply)
          + key      = "0x100006"
          + one      = "00:00:00:00:00:01/None"
          + priority = 100
          + two      = "00:00:00:00:00:02/None"
          + type     = "HostToHostIntent"
        }
      + last_updated = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + h1-to-h2_intent = {
      + id           = (known after apply)
      + intent       = {
          + appid    = "org.onosproject.cli"
          + id       = (known after apply)
          + key      = "0x100006"
          + one      = "00:00:00:00:00:01/None"
          + priority = 100
          + two      = "00:00:00:00:00:02/None"
          + type     = "HostToHostIntent"
        }
      + last_updated = (known after apply)
    }
```

#### Hosts
Hosts can be pulled from onos as a data source for use in configuration for intents.

Configuration:
```hcl
data "onos_hosts" "mininet" {}

output "mininet_hosts" {
  value = data.onos_hosts.mininet
}
```

Output:
```shell
$ terraform plan
data.onos_hosts.mininet: Reading...
data.onos_hosts.mininet: Read complete after 0s [id=placeholder]

Changes to Outputs:
  + mininet_hosts = {
      + hosts = [
          + {
              + configured  = false
              + id          = "00:00:00:00:00:03/None"
              + innervlan   = "None"
              + ipaddresses = [
                  + "10.0.0.3",
                ]
              + locations   = [
                  + {
                      + elementid = "of:0000000000000003"
                      + port      = "1"
                    },
                ]
              + mac         = "00:00:00:00:00:03"
              + outertpid   = "0x0000"
              + suspended   = false
              + vlan        = "None"
            },
            # Truncated
```

#### Flows
Flows can be pulled from onos as a data source. The flows data source was added early during development. There doesn't currently seem to be a need for use of flows in other configurations, but it has been left for future use cases and as a reference for how to deal with various types and schemas.

Configuration:
```hcl
data "onos_flows" "mininet" {}

output "mininet_flows" {
  value = data.onos_flows.mininet
}
```

Output:
```shell
$ terraform plan

data.onos_flows.mininet: Reading...
data.onos_flows.mininet: Read complete after 0s

Changes to Outputs:
  + mininet_flows = {
      + flows = [
          + {
              + appid       = "org.onosproject.core"
              + bytes       = 187094
              + deviceid    = "of:0000000000000003"
              + groupid     = 0
              + id          = "281476661728682"
              + ispermanent = true
              + lastseen    = 1700167313863
              + life        = 4175
              + livetype    = "UNKNOWN"
              + packets     = 1346
              + priority    = 40000
              + selector    = {
                  + criteria = [
                      + {
                          + ethtype = "0x88cc"
                          + mac     = ""
                          + port    = 0
                          + type    = "ETH_TYPE"
                        },
                    ]
                }
              + state       = "ADDED"
              + tableid     = 0
              + tablename   = "0"
              + timeout     = 0
              + treatment   = {
                  + cleardeferred = true
                  + deferred      = []
                  + instructions  = [
                      + {
                          + port = "CONTROLLER"
                          + type = "OUTPUT"
                        },
                    ]
                }
            },
            # Truncated
```


## Example Using Docker Containers running on Linux (Ubuntu 22.04.3 LTS)
These examples require a current version of [go](https://go.dev/doc/install) and [docker](https://docs.docker.com/engine/install/ubuntu/).

#### Docker Environment Setup
Navigate to the examples directory and run:

```shell
sudo docker compose up
```

This runs onos and mininet in docker containers and links them.  

From a new terminal, paste in the following curl commands to activate openflow and fwd in ONOS:

```shell
curl --request POST \
--url http://localhost:8181/onos/v1/applications/org.onosproject.fwd/active \
--header 'Accept: application/json' \
--header 'Authorization: Basic b25vczpyb2Nrcw=='

curl --request POST \
--url http://localhost:8181/onos/v1/applications/org.onosproject.openflow/active \
--header 'Accept: application/json' \
--header 'Authorization: Basic b25vczpyb2Nrcw=='
```

The calls should return output similar to this, indicating the activation was successful:

```
{"name":"org.onosproject.fwd","id":78,"version":"2.7.1.SNAPSHOT","category":"Traffic Engineering","description":"Provisions traffic between end-stations using hop-by-hop flow programming by intercepting packets for which there are currently no matching flow objectives on the data plane.","readme":"Provisions traffic between end-stations using hop-by-hop flow programming by intercepting packets for which there are currently no matching flow objectives on the data plane. The paths paved in this manner are short-lived, i.e. they expire a few seconds after the flow on whose behalf they were programmed stops.\\n\\nThe application relies on the ONOS path service to compute the shortest paths. In the event of negative topology events (link loss, device disconnect, etc.), the application will proactively invalidate any paths that it had programmed to lead through the resources that are no longer available.","origin":"ONOS Community","url":"http://onosproject.org","featuresRepo":"mvn:org.onosproject/onos-apps-fwd/2.7.1-SNAPSHOT/xml/features","state":"ACTIVE","features":["onos-apps-fwd"],"permissions":[],"requiredApps":[]}

{"name":"org.onosproject.openflow","id":46,"version":"2.7.1.SNAPSHOT","category":"Provider","description":"Suite of the OpenFlow base providers bundled together with ARP\\/NDP host location provider and LLDP link provider.","readme":"Suite of the OpenFlow base providers bundled together with ARP\\/NDP host location provider and LLDP link provider.","origin":"ONOS Community","url":"http://onosproject.org","featuresRepo":"mvn:org.onosproject/onos-providers-openflow-app/2.7.1-SNAPSHOT/xml/features","state":"ACTIVE","features":["onos-providers-openflow-app"],"permissions":[],"requiredApps":["org.onosproject.hostprovider","org.onosproject.lldpprovider","org.onosproject.openflow-base"]}
```

Launch a terminal in the Mininet container:

```shell
sudo docker exec -it mininet-terraform /bin/bash
```

Start a basic mininet topology with onos as the SDN controller:

```shell
mn --topo tree,2,2 --mac --switch ovs,protocols=OpenFlow14 --controller remote,ip=onos-terraform
```

Mininet should successfully connect to the ONOS controller, start the switches and hosts, and display the mininet CLI prompt:

```shell
root@c0b74ca1c8a7:~# mn --topo tree,2,2 --mac --switch ovs,protocols=OpenFlow14   --controller remote,ip=onos-terraform
*** Error setting resource limits. Mininet's performance may be affected.
*** Creating network
*** Adding controller
Connecting to remote controller at onos-compose:6653
*** Adding hosts:
h1 h2 h3 h4 
*** Adding switches:
s1 s2 s3 
*** Adding links:
(s1, s2) (s1, s3) (s2, h1) (s2, h2) (s3, h3) (s3, h4) 
*** Configuring hosts
h1 h2 h3 h4 
*** Starting controller
c0 
*** Starting 3 switches
s1 s2 s3 ...
*** Starting CLI:
mininet> 
```

From the mininet CLI, run pingall and confirm that all hosts can communicate:
```shell
mininet> pingall
*** Ping: testing ping reachability
h1 -> h2 h3 h4 
h2 -> h1 h3 h4 
h3 -> h1 h2 h4 
h4 -> h1 h2 h3 
*** Results: 0% dropped (12/12 received)
```

Switch to an ubuntu terminal (not the mininet container) and paste in the following curl command to deactivate onos fwd:
```shell
curl -X DELETE --header 'Accept: application/json' 'http://localhost:8181/onos/v1/applications/org.onosproject.fwd/active'
```

This disables reactive forwarding in the onos controller, causing to to behave more like a firewall than a router.

From the mininet CLI, run pingall again and confirm that the traffic is now being blocked:
```shell
mininet> pingall
*** Ping: testing ping reachability
h1 -> X X X 
h2 -> X X X 
h3 -> X X X 
h4 -> X X X 
*** Results: 100% dropped (0/12 received)
```
### Provider Demonstration

#### Create an Intent

Navigate to the examples/intent directory, create a file named main.tf and paste in the following configuration:

```hcl
terraform {
  required_providers {
    onos = {
      source = "hashicorp.com/cjtnkns/onos"
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
```

Run ``terraform plan`` from the intent directory. Terraform will display the changes to be made.

```shell
$ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # onos_intent.h1-to-h2 will be created
  + resource "onos_intent" "h1-to-h2" {
      + id           = (known after apply)
      + intent       = {
          + appid    = "org.onosproject.cli"
          + id       = (known after apply)
          + key      = "0x100006"
          + one      = "00:00:00:00:00:01/None"
          + priority = 100
          + two      = "00:00:00:00:00:02/None"
          + type     = "HostToHostIntent"
        }
      + last_updated = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + h1-to-h2_intent = {
      + id           = (known after apply)
      + intent       = {
          + appid    = "org.onosproject.cli"
          + id       = (known after apply)
          + key      = "0x100006"
          + one      = "00:00:00:00:00:01/None"
          + priority = 100
          + two      = "00:00:00:00:00:02/None"
          + type     = "HostToHostIntent"
        }
      + last_updated = (known after apply)
    }

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
```

Run: `terraform apply -auto-approve` to apply the configuration:

```shell
$ terraform apply -auto-approve

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # onos_intent.h1-to-h2 will be created
  + resource "onos_intent" "h1-to-h2" {
      + id           = (known after apply)
      + intent       = {
          + appid    = "org.onosproject.cli"
          + id       = (known after apply)
          + key      = "0x100006"
          + one      = "00:00:00:00:00:01/None"
          + priority = 100
          + two      = "00:00:00:00:00:02/None"
          + type     = "HostToHostIntent"
        }
      + last_updated = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + h1-to-h2_intent = {
      + id           = (known after apply)
      + intent       = {
          + appid    = "org.onosproject.cli"
          + id       = (known after apply)
          + key      = "0x100006"
          + one      = "00:00:00:00:00:01/None"
          + priority = 100
          + two      = "00:00:00:00:00:02/None"
          + type     = "HostToHostIntent"
        }
      + last_updated = (known after apply)
    }
onos_intent.h1-to-h2: Creating...
onos_intent.h1-to-h2: Creation complete after 0s [id=0x6]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

h1-to-h2_intent = {
  "id" = "0x6"
  "intent" = {
    "appid" = "org.onosproject.cli"
    "id" = "0x6"
    "key" = "0x100006"
    "one" = "00:00:00:00:00:01/None"
    "priority" = 100
    "two" = "00:00:00:00:00:02/None"
    "type" = "HostToHostIntent"
  }
  "last_updated" = "Wednesday, 03-Jan-24 13:56:21 EST"
}
```

Notice in the output that the id field has been computed and updated from the ONOS response.

Paste the following command in the terminal to retrieve the intents from ONOS and confirm the intent was created successfully:

```shell
curl --request GET --url http://127.0.0.1:8181/onos/v1/intents --header 'Accept: application/json' --header 'Authorization: Basic b25vczpyb2Nrcw=='
```

```json
{"intents":[{"type":"HostToHostIntent","id":"0x6","key":"0x100006","appId":"org.onosproject.cli","resources":["00:00:00:00:00:01/None","00:00:00:00:00:02/None"],"state":"INSTALLED"}]}
```

**Note: If the "One" and "Two" strings do not match any hosts from your onos environmnet, the intent may fail to install. Edit the one and two strings in main.tf to match two of the hosts in your environment.**

From the mininet CLI, test connectivity again and confirm that the traffic between h1 and h2 is now allowed:

```shell
mininet> h1 ping -c 4 h2
PING 10.0.0.2 (10.0.0.2) 56(84) bytes of data.
64 bytes from 10.0.0.2: icmp_seq=1 ttl=64 time=0.556 ms
64 bytes from 10.0.0.2: icmp_seq=2 ttl=64 time=0.038 ms
64 bytes from 10.0.0.2: icmp_seq=3 ttl=64 time=0.036 ms
64 bytes from 10.0.0.2: icmp_seq=4 ttl=64 time=0.038 ms

--- 10.0.0.2 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3063ms
rtt min/avg/max/mdev = 0.036/0.167/0.556/0.224 ms
```

Confirm that traffic not allowed by the intent is still blocked:
```shell
mininet> h1 ping -c 4 h3 
PING 10.0.0.3 (10.0.0.3) 56(84) bytes of data.
From 10.0.0.1 icmp_seq=1 Destination Host Unreachable
From 10.0.0.1 icmp_seq=2 Destination Host Unreachable
From 10.0.0.1 icmp_seq=3 Destination Host Unreachable
From 10.0.0.1 icmp_seq=4 Destination Host Unreachable

--- 10.0.0.3 ping statistics ---
4 packets transmitted, 0 received, +4 errors, 100% packet loss, time 3078ms
pipe 4
```

#### Update an Intent
Replace the h1-to-h2 resource configuration in main.tf with this:

```hcl
resource "onos_intent" "h1-to-h2" {
  intent = {
    appid    = "org.onosproject.cli"
    key      = "0x100006"
    type     = "HostToHostIntent"
    priority = 100
    one      = "00:00:00:00:00:01/None"
    two      = "00:00:00:00:00:03/None"
  }
}
```

Run `terraform plan`. Terraform will display the changes to be made.


```shell
$ terraform plan

onos_intent.h1-to-h2: Refreshing state... [id=0x6]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # onos_intent.h1-to-h2 will be updated in-place
  ~ resource "onos_intent" "h1-to-h2" {
      ~ id           = "0x6" -> (known after apply)
      ~ intent       = {
          ~ id       = "0x6" -> (known after apply)
          ~ two      = "00:00:00:00:00:02/None" -> "00:00:00:00:00:03/None"
            # (5 unchanged attributes hidden)
        }
      ~ last_updated = "Wednesday, 03-Jan-24 13:56:21 EST" -> (known after apply)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ h1-to-h2_intent = {
      ~ id           = "0x6" -> (known after apply)
      ~ intent       = {
          ~ id       = "0x6" -> (known after apply)
          ~ two      = "00:00:00:00:00:02/None" -> "00:00:00:00:00:03/None"
            # (5 unchanged attributes hidden)
        }
      ~ last_updated = "Wednesday, 03-Jan-24 13:56:21 EST" -> (known after apply)
    }

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
```

The "two" field displays an in-place update: `~ two      = "00:00:00:00:00:02/None" -> "00:00:00:00:00:03/None"`

Run: `terraform apply -auto-approve` to apply the configuration:

```shell
$ terraform apply -auto-approve

onos_intent.h1-to-h2: Refreshing state... [id=0x6]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # onos_intent.h1-to-h2 will be updated in-place
  ~ resource "onos_intent" "h1-to-h2" {
      ~ id           = "0x6" -> (known after apply)
      ~ intent       = {
          ~ id       = "0x6" -> (known after apply)
          ~ two      = "00:00:00:00:00:02/None" -> "00:00:00:00:00:03/None"
            # (5 unchanged attributes hidden)
        }
      ~ last_updated = "Wednesday, 03-Jan-24 13:56:21 EST" -> (known after apply)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ h1-to-h2_intent = {
      ~ id           = "0x6" -> (known after apply)
      ~ intent       = {
          ~ id       = "0x6" -> (known after apply)
          ~ two      = "00:00:00:00:00:02/None" -> "00:00:00:00:00:03/None"
            # (5 unchanged attributes hidden)
        }
      ~ last_updated = "Wednesday, 03-Jan-24 13:56:21 EST" -> (known after apply)
    }
onos_intent.h1-to-h2: Modifying... [id=0x6]
onos_intent.h1-to-h2: Modifications complete after 0s [id=0xb]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

h1-to-h2_intent = {
  "id" = "0xb"
  "intent" = {
    "appid" = "org.onosproject.cli"
    "id" = "0xb"
    "key" = "0x100006"
    "one" = "00:00:00:00:00:01/None"
    "priority" = 100
    "two" = "00:00:00:00:00:03/None"
    "type" = "HostToHostIntent"
  }
  "last_updated" = "Wednesday, 03-Jan-24 14:07:33 EST"
}
```

The state in the output now reflects the updated configuration.

Paste the following command in the terminal to retrieve the intents from ONOS and confirm the intent was created successfully:

```shell
curl --request GET --url http://127.0.0.1:8181/onos/v1/intents --header 'Accept: application/json' --header 'Authorization: Basic b25vczpyb2Nrcw=='
```

```json
{"intents":[{"type":"HostToHostIntent","id":"0xb","key":"0x100006","appId":"org.onosproject.cli","resources":["00:00:00:00:00:01/None","00:00:00:00:00:03/None"],"state":"INSTALLED"}]}
```

Confirm that h1 to h2 is now blocked as expected and traffic from h1 to h3 is allowed by the intent:
```shell
mininet> h1 ping -c 4 h2
PING 10.0.0.2 (10.0.0.2) 56(84) bytes of data.


--- 10.0.0.2 ping statistics ---
4 packets transmitted, 0 received, 100% packet loss, time 3782ms

mininet> h1 ping -c 4 h3
PING 10.0.0.3 (10.0.0.3) 56(84) bytes of data.
64 bytes from 10.0.0.3: icmp_seq=1 ttl=64 time=0.334 ms
64 bytes from 10.0.0.3: icmp_seq=2 ttl=64 time=0.043 ms
64 bytes from 10.0.0.3: icmp_seq=3 ttl=64 time=0.048 ms
64 bytes from 10.0.0.3: icmp_seq=4 ttl=64 time=0.068 ms

--- 10.0.0.3 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3058ms
rtt min/avg/max/mdev = 0.043/0.123/0.334/0.122 ms
```
#### Import an Intent
Paste in the following curl command to add an intent to ONOS:
```shell
curl --request POST \
  --url http://localhost:8181/onos/v1/intents \
  --header 'Accept: application/json' \
  --header 'Authorization: Basic b25vczpyb2Nrcw==' \
  --header 'Content-Type: application/json' \
  --data '{
			"type": "HostToHostIntent",
			"appId": "org.onosproject.cli",
	        "key": "0x100007",
			"one": "00:00:00:00:00:02/None",
			"two": "00:00:00:00:00:03/None"
		}'
```


Get the intents and confirm that the new intent was created.
```shell
curl --request GET --url http://127.0.0.1:8181/onos/v1/intents --header 'Accept: application/json' --header 'Authorization: Basic b25vczpyb2Nrcw=='
```

Both intents should no be shown.
```json
{"intents":[{"type":"HostToHostIntent","id":"0x10","key":"0x100007","appId":"org.onosproject.cli","resources":["00:00:00:00:00:02/None","00:00:00:00:00:03/None"],"state":"INSTALLED"},{"type":"HostToHostIntent","id":"0xb","key":"0x100006","appId":"org.onosproject.cli","resources":["00:00:00:00:00:01/None","00:00:00:00:00:03/None"],"state":"INSTALLED"}]}
```

We will now import that intent into the Terraform state so that it can be managed by Terraform. 

Add the following lines to the main.tf file:

```hcl
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

```

Run `terraform state list` and confirm that currently the only state file is "h1-to-h2":
```shell
$ terraform state list
onos_intent.h1-to-h2
```

Run `terraform import onos_intent.h2-to-h3 "org.onosproject.cli,0x100007"`:

The first value after import matches the output value from the main.tf configuration. This determines the name terraform will use to manage the state of the intent. The next field is a comma separated string containing the appid and key of the intent. These values must match the intent that exists in ONOS since they are used to lookup and import the intent. 

```shell
 terraform import onos_intent.h2-to-h3 "org.onosproject.cli,0x100007terraform import onos_intent.h2-to-h3 "org.onosproject.cli,0x100007"
onos_intent.h2-to-h3: Importing from ID "org.onosproject.cli,0x100007"...
onos_intent.h2-to-h3: Import prepared!
  Prepared onos_intent for import
onos_intent.h2-to-h3: Refreshing state...

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Run `terraform state list` and confirm that both "h1-to-h2" and "h2-to-h3" are now shown.:

```shell
$ terraform state list
onos_intent.h1-to-h2
onos_intent.h2-to-h3
```

Test connectivity from h2 to h3 to confirm the new intent works as expected.

```shell
mininet> h2 ping -c 4 h3
PING 10.0.0.3 (10.0.0.3) 56(84) bytes of data.
64 bytes from 10.0.0.3: icmp_seq=1 ttl=64 time=0.192 ms
64 bytes from 10.0.0.3: icmp_seq=2 ttl=64 time=0.058 ms
64 bytes from 10.0.0.3: icmp_seq=3 ttl=64 time=0.059 ms
64 bytes from 10.0.0.3: icmp_seq=4 ttl=64 time=0.041 ms

--- 10.0.0.3 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3071ms
rtt min/avg/max/mdev = 0.041/0.087/0.192/0.060 ms
```

A `terraform plan` should show the infrastructure matches the configuration:
```shell
$ terraform plan

onos_intent.h1-to-h2: Refreshing state... [id=0xb]
onos_intent.h2-to-h3: Refreshing state... [id=0x10]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

#### Delete an intent
Delete the "h1-to-h2" resource and output blocks:

```hcl
resource "onos_intent" "h1-to-h2" {
  intent = {
    appid    = "org.onosproject.cli"
    key      = "0x100006"
    type     = "HostToHostIntent"
    priority = 100
    one      = "00:00:00:00:00:01/None"
    two      = "00:00:00:00:00:03/None"
  }
}

output "h1-to-h2_intent" {
  value = onos_intent.h1-to-h2
}
```

Run `terraform plan`, note that the resource will be destroyed.:

```shell
$ terraform plan

onos_intent.h2-to-h3: Refreshing state... [id=0x10]
onos_intent.h1-to-h2: Refreshing state... [id=0xb]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # onos_intent.h1-to-h2 will be destroyed
  # (because onos_intent.h1-to-h2 is not in configuration)
  - resource "onos_intent" "h1-to-h2" {
      - id           = "0xb" -> null
      - intent       = {
          - appid    = "org.onosproject.cli" -> null
          - id       = "0xb" -> null
          - key      = "0x100006" -> null
          - one      = "00:00:00:00:00:01/None" -> null
          - priority = 100 -> null
          - two      = "00:00:00:00:00:03/None" -> null
          - type     = "HostToHostIntent" -> null
        } -> null
      - last_updated = "Wednesday, 03-Jan-24 14:07:33 EST" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - h1-to-h2_intent = {
      - id           = "0xb"
      - intent       = {
          - appid    = "org.onosproject.cli"
          - id       = "0xb"
          - key      = "0x100006"
          - one      = "00:00:00:00:00:01/None"
          - priority = 100
          - two      = "00:00:00:00:00:03/None"
          - type     = "HostToHostIntent"
        }
      - last_updated = "Wednesday, 03-Jan-24 14:07:33 EST"
    } -> null

────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
```

Run: `terraform apply -auto-approve` to apply the changes:

```shell
$ terraform apply -auto-approve

onos_intent.h1-to-h2: Refreshing state... [id=0xb]
onos_intent.h2-to-h3: Refreshing state... [id=0x10]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # onos_intent.h1-to-h2 will be destroyed
  # (because onos_intent.h1-to-h2 is not in configuration)
  - resource "onos_intent" "h1-to-h2" {
      - id           = "0xb" -> null
      - intent       = {
          - appid    = "org.onosproject.cli" -> null
          - id       = "0xb" -> null
          - key      = "0x100006" -> null
          - one      = "00:00:00:00:00:01/None" -> null
          - priority = 100 -> null
          - two      = "00:00:00:00:00:03/None" -> null
          - type     = "HostToHostIntent" -> null
        } -> null
      - last_updated = "Wednesday, 03-Jan-24 14:07:33 EST" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - h1-to-h2_intent = {
      - id           = "0xb"
      - intent       = {
          - appid    = "org.onosproject.cli"
          - id       = "0xb"
          - key      = "0x100006"
          - one      = "00:00:00:00:00:01/None"
          - priority = 100
          - two      = "00:00:00:00:00:03/None"
          - type     = "HostToHostIntent"
        }
      - last_updated = "Wednesday, 03-Jan-24 14:07:33 EST"
    } -> null
onos_intent.h1-to-h2: Destroying... [id=0xb]
onos_intent.h1-to-h2: Destruction complete after 0s

Apply complete! Resources: 0 added, 0 changed, 1 destroyed.

Outputs:

h2-to-h3_intent = {
  "id" = "0x10"
  "intent" = {
    "appid" = "org.onosproject.cli"
    "id" = "0x10"
    "key" = "0x100007"
    "one" = "00:00:00:00:00:02/None"
    "priority" = 100
    "two" = "00:00:00:00:00:03/None"
    "type" = "HostToHostIntent"
  }
  "last_updated" = tostring(null)
}
```

Paste the following command in the terminal to retrieve the intents from ONOS and confirm the intent was created successfully:

```shell
curl --request GET --url http://127.0.0.1:8181/onos/v1/intents --header 'Accept: application/json' --header 'Authorization: Basic b25vczpyb2Nrcw=='
```

The original intent is gone and only the h2-h3 one remains.
```json
{"intents":[{"type":"HostToHostIntent","id":"0x10","key":"0x100007","appId":"org.onosproject.cli","resources":["00:00:00:00:00:02/None","00:00:00:00:00:03/None"],"state":"INSTALLED"}]}
```

Confirm pings work as expected:
```shell
mininet> h2 ping -c 4 h3
PING 10.0.0.3 (10.0.0.3) 56(84) bytes of data.
64 bytes from 10.0.0.3: icmp_seq=1 ttl=64 time=711 ms
64 bytes from 10.0.0.3: icmp_seq=2 ttl=64 time=0.048 ms
64 bytes from 10.0.0.3: icmp_seq=3 ttl=64 time=0.044 ms
64 bytes from 10.0.0.3: icmp_seq=4 ttl=64 time=0.041 ms

--- 10.0.0.3 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3036ms
rtt min/avg/max/mdev = 0.041/177.885/711.408/308.029 ms

mininet> h1 ping -c 4 h3
PING 10.0.0.3 (10.0.0.3) 56(84) bytes of data.
From 10.0.0.1 icmp_seq=1 Destination Host Unreachable
From 10.0.0.1 icmp_seq=2 Destination Host Unreachable
From 10.0.0.1 icmp_seq=3 Destination Host Unreachable
From 10.0.0.1 icmp_seq=4 Destination Host Unreachable

--- 10.0.0.3 ping statistics ---
4 packets transmitted, 0 received, +4 errors, 100% packet loss, time 3073ms
pipe 4
```

## Testing
Automated acceptance testing has been implemented in accordance with Terraform's best practices for providers.

```shell
$ TF_ACC=1 go test -count=1 -v
=== RUN   TestAccHostsDataSource
--- PASS: TestAccHostsDataSource (3.59s)
=== RUN   TestAccIntentResource
--- PASS: TestAccIntentResource (4.74s)
PASS
ok      terraform-provider-onos/internal/provider       8.354s
```

## Limitations
Care must be taken when editing the configration for exisitng intents. Updating the "one" and "two" fields has been thouroughly tested and is stable. **The appid and key should never be edited**. Ideally these values would be computed instead of configured. Unfortunately, when a new intent is created in ONOS without specifying the key, a random key is used and that key is not returned. Without the key, there is no way to look up the intent for future operations or even confirm that it was created successfully. 

In the future, the ONOS API may be updated to return the intent details (including the key) when an intent is created. A possible work around is to generate a key using a hash or some other method in the ONOS API Client library and inject that key into the CreateIntent() API call so that it does not have to be specified in the terraform configuration, but this is not an ideal solution and has not bee necessary yet. 

```hcl
resource "onos_intent" "h1-to-h2" {
  intent = {
    appid    = "org.onosproject.cli"    #Do not edit
    key      = "0x100006"               #Do not edit
    type     = "HostToHostIntent"
    priority = 100
    one      = "00:00:00:00:00:01/None"
    two      = "00:00:00:00:00:02/None"
  }
}
```
