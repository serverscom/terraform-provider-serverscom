package serverscom

import (
	"context"
	"fmt"
	"log"
	"regexp"
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
		return fmt.Errorf("error getting client for sweeping dedicated servers: %s", err)
	}

	ctx := context.TODO()

	hosts, err := client.Hosts.Collection().Collect(ctx)
	if err != nil {
		return fmt.Errorf("error getting list of hosts for sweeping dedicated servers: %s", err)
	}

	for _, host := range hosts {
		_, err := client.Hosts.ScheduleReleaseForDedicatedServer(ctx, host.ID, scgo.ScheduleReleaseInput{})
		if err != nil {
			return fmt.Errorf("can't schedule release for dedicated server (%s): %s", host.ID, err)
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
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no dedicated server ID is set")
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
				return fmt.Errorf("error retrieving dedicated server: %s", err)
			}
		}

		if dedicatedServer.ScheduledRelease == nil {
			return fmt.Errorf("Dedicated server (%s) has not been scheduled to release", rs.Primary.ID)
		}
	}

	return nil
}

func TestDedicatedServerWillReinstall(t *testing.T) {
	cases := []struct {
		name       string
		oldTrigger string
		newTrigger string
		want       bool
	}{
		{"unchanged none", "none", "none", false},
		{"first adoption keeps none", "none", "none", false},
		{"empty to none is a no-op", "", "none", false},
		{"reset back to none", "1", "none", false},
		{"unchanged value", "1", "1", false},
		{"bump from none", "none", "1", true},
		{"bump between values", "1", "2", true},
		{"explicit set from empty", "", "1", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := dedicatedServerWillReinstall(tc.oldTrigger, tc.newTrigger)
			if got != tc.want {
				t.Fatalf("dedicatedServerWillReinstall(%q, %q) = %v, want %v",
					tc.oldTrigger, tc.newTrigger, got, tc.want)
			}
		})
	}
}

func TestAccServerscomDedicatedServer_Reinstall(t *testing.T) {
	var dedicatedServer scgo.DedicatedServer
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckDedicatedServerDestroy,
		Steps: []resource.TestStep{
			{
				// Create with reinstall_trigger at its default "none".
				Config: testAccServerscomCheckDedicatedServerConfig_reinstall(rInt, "Ubuntu 22.04-server x86_64", "none"),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDedicatedServerExists("serverscom_dedicated_server.node", &dedicatedServer),
					resource.TestCheckResourceAttr(
						"serverscom_dedicated_server.node", "reinstall_trigger", "none"),
				),
			},
			{
				// Changing operating_system without bumping reinstall_trigger must fail the plan.
				Config:      testAccServerscomCheckDedicatedServerConfig_reinstall(rInt, "Ubuntu 24.04-server x86_64", "none"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("requires an OS reinstall"),
			},
			{
				// Bumping reinstall_trigger acknowledges and performs the reinstall.
				Config: testAccServerscomCheckDedicatedServerConfig_reinstall(rInt, "Ubuntu 24.04-server x86_64", "1"),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDedicatedServerExists("serverscom_dedicated_server.node", &dedicatedServer),
					resource.TestCheckResourceAttr(
						"serverscom_dedicated_server.node", "operating_system", "Ubuntu 24.04-server x86_64"),
					resource.TestCheckResourceAttr(
						"serverscom_dedicated_server.node", "reinstall_trigger", "1"),
					resource.TestCheckResourceAttr(
						"serverscom_dedicated_server.node", "operational_status", "normal"),
				),
			},
		},
	})
}

func testAccServerscomCheckDedicatedServerConfig_reinstall(rInt int, operatingSystem, reinstallTrigger string) string {
	return fmt.Sprintf(`
resource "serverscom_dedicated_server" "node" {
	bandwidth            = "19.1 TB"
	hostname             = "node-%d"
	location             = "SJC1"
	operating_system     = "%s"
	reinstall_trigger    = "%s"
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
			target = "swap"
			size = 4096
			fill = false
		}
	}
}
`, rInt, operatingSystem, reinstallTrigger)
}
