package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "onos_hosts" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of hosts returned
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.#", "4"),
					// Verify the first host to ensure all attributes are set
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.id", "00:00:00:00:00:03/None"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.mac", "00:00:00:00:00:03"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.vlan", "None"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.innervlan", "None"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.configured", "false"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.suspended", "false"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.ipaddresses.#", "1"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.ipaddresses.0", "10.0.0.3"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.locations.#", "1"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.locations.0.elementid", "of:0000000000000003"),
					resource.TestCheckResourceAttr("data.onos_hosts.test", "hosts.0.locations.0.port", "1"),

					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.onos_hosts.test", "id", "placeholder"),
				),
			},
		},
	})
}
