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
	resource.AddTestSweepers("serverscom_sbm_server", &resource.Sweeper{
		Name: "serverscom_sbm_server",
		F:    testSweepSBMServers,
	})
}

func testSweepSBMServers(region string) error {
	log.Printf("[DEBUG] Sweeping SBM servers")
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping SBM servers: %s", err)
	}

	ctx := context.TODO()

	servers, err := client.Hosts.Collection().Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error getting list of SBM servers: %s", err)
	}

	for _, server := range servers {
		_, err := client.Hosts.ReleaseSBMServer(ctx, server.ID)
		if err != nil {
			return fmt.Errorf("Can't release SBM server (%s): %s", server.ID, err)
		}
	}

	return nil
}

func TestAccServerscomSBMServer_Basic(t *testing.T) {
	var sbmServer scgo.SBMServer
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckSBMServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomSBMServerConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckSBMServerExists("serverscom_sbm_server.node", &sbmServer),
					resource.TestCheckResourceAttr(
						"serverscom_sbm_server.node", "hostname", fmt.Sprintf("node-%d", rInt)),
					resource.TestCheckResourceAttr(
						"serverscom_sbm_server.node", "operating_system", "Ubuntu 18.04-server x86_64"),
				),
			},
		},
	})
}

func testAccServerscomSBMServerConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "serverscom_sbm_server" "node" {
	hostname         = "node-%d"
	location         = "SJC1"
	flavor           = "SBM-01"
	operating_system = "Ubuntu 18.04-server x86_64"
}
`, rInt)
}

func testAccServerscomCheckSBMServerExists(n string, sbmServer *scgo.SBMServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SBM server ID is set")
		}

		client := testAccProvider.Meta().(*scgo.Client)

		currentSBMServer, err := client.Hosts.GetSBMServer(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		*sbmServer = *currentSBMServer
		return nil
	}
}

func testAccServerscomCheckSBMServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*scgo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "serverscom_sbm_server" {
			continue
		}

		sbmServer, err := client.Hosts.GetSBMServer(context.Background(), rs.Primary.ID)
		if err != nil {
			switch err.(type) {
			case *scgo.NotFoundError:
				return nil
			default:
				return fmt.Errorf("Error retrieving SBM server: %s", err)
			}
		}

		if sbmServer.Status != "released" {
			return fmt.Errorf("SBM server (%s) has not been released", rs.Primary.ID)
		}
	}

	return nil
}
