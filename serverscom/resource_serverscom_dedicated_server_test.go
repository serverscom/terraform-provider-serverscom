package serverscom

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func init() {
	resource.AddTestSweepers("serverscom_dedicated_server", &resource.Sweeper{
		Name: "serverscom_dedicated_server",
		F:    testSweepDededicatedServers,
	})
}

func testSweepDededicatedServers(region string) error {
	log.Printf("[DEBUG] Sweeping dedicated servers")
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping dedicated servers: %s", err)
	}

	ctx := context.TODO()

	hosts, err := client.Hosts.Collection().Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error getting list of hosts for sweeping dedicated servers: %s", err)
	}

	for _, host := range hosts {
		_, err := client.Hosts.ScheduleReleaseForDedicatedServer(ctx, host.ID, scgo.ScheduleReleaseInput{})
		if err != nil {
			return fmt.Errorf("Can't schedule release for dedicated server (%s): %s", host.ID, err)
		}
	}

	return nil
}

func TestAccServerscomDedicatedServer_Basic(t *testing.T) {
	var dedicatedServer scgo.DedicatedServer
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckDedicatedServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomCheckDedicatedServerConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDedicatedServerExists("serverscom_dedicated_server.node", &dedicatedServer),
					resource.TestCheckResourceAttr(
						"serverscom_dedicated_server.node", "name", fmt.Sprintf("node-%d", rInt)),
					resource.TestCheckResourceAttr(
						"serverscom_dedicated_server.node", "operating_system", "Ubuntu 16.04-server x86_64"),
				),
			},
		},
	})
}

func testAccServerscomCheckDedicatedServerConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "serverscom_dedicated_server" "node" {
	bandwidth            = "19.1 TB"
	hostname             = "node-%d"
	location             = "SJC1"
	operating_system     = "Ubuntu 16.04-server x86_64"
	private_uplink       = "Private 10 Gbps with redundancy"
	public_uplink        = "Public 10 Gbps with redundancy"
	ram_size             = 32
	server_model         = "Dell R440 / 2xIntel Xeon Silver-4114 / 32 GB RAM / 1x480 GB SSD"

	slot {
		drive_model = "480 GB SSD SATA"
		position    = 0
	}

	layout {
		slot_positions = [0]

		partition {
			target = "/"
			size = 10240
			fill = false
			fs = "ext4"
		}

		partition {
			target = "/home"
			size = 1
			fill = true
			fs = "ext4"
		}

		partition {
			target = "swap"
			size = 4096
			fill = false
		}
	}
`, rInt)
}

func testAccServerscomCheckDedicatedServerExists(n string, dedicatedServer *scgo.DedicatedServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No dedicated server ID is set")
		}

		client := testAccProvider.Meta().(*scgo.Client)

		currentDedicatedServer, err := client.Hosts.GetDedicatedServer(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		*dedicatedServer = *currentDedicatedServer
		return nil
	}
}

func testAccServerscomCheckDedicatedServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*scgo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "serverscom_dedicated_server" {
			continue
		}

		dedicatedServer, err := client.Hosts.GetDedicatedServer(context.Background(), rs.Primary.ID)
		if err != nil {
			switch err.(type) {
			case *scgo.NotFoundError:
				return nil
			default:
				return fmt.Errorf("Error retrieving dedicated server: %s", err)
			}
		}

		if dedicatedServer.ScheduledRelease == nil {
			return fmt.Errorf("Dedicated server (%s) has not been scheduled to release", rs.Primary.ID)
		}
	}

	return nil
}
