package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "onos_intent" "test" {
					intent = {
					  appid    = "org.onosproject.cli"
					  key      = "0x99999"
					  type     = "HostToHostIntent"
					  priority = 100
					  one      = "00:00:00:00:00:99/None"
					  two      = "00:00:00:00:00:00/None"
					}
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					//resource.TestCheckResourceAttr("onos_intent.test", "items.#", "1"),
					// Verify first intent item
					resource.TestCheckResourceAttr("onos_intent.test", "intent.appid", "org.onosproject.cli"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.key", "0x99999"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.priority", "100"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.type", "HostToHostIntent"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.one", "00:00:00:00:00:99/None"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.two", "00:00:00:00:00:00/None"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("onos_intent.test", "id"),
					resource.TestCheckResourceAttrSet("onos_intent.test", "intent.id"),
					resource.TestCheckResourceAttrSet("onos_intent.test", "last_updated"),
				),
			},
			// ImportState testing
			/*
				{
					ResourceName:      "onos_intent.test",
					ImportState:       true,
					ImportStateVerify: true,
					// The last_updated attribute does not exist in the onos
					// API, therefore there is no value for it during import.
					ImportStateVerifyIgnore: []string{"last_updated"},
				},
				// Update and Read testing
			*/
			{
				Config: providerConfig + `
				resource "onos_intent" "test" {
					intent = {
					  appid    = "org.onosproject.cli"
					  key      = "0x99999"
					  type     = "HostToHostIntent"
					  priority = 100
					  one      = "00:00:00:00:00:98/None"
					  two      = "00:00:00:00:00:01/None"
					}
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify intent updated
					resource.TestCheckResourceAttr("onos_intent.test", "intent.appid", "org.onosproject.cli"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.key", "0x99999"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.priority", "100"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.type", "HostToHostIntent"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.one", "00:00:00:00:00:98/None"),
					resource.TestCheckResourceAttr("onos_intent.test", "intent.two", "00:00:00:00:00:01/None"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
