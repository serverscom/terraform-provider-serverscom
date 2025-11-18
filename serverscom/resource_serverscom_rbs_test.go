package serverscom

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	// set actual location and flavor before runing acceptance tests
	testLocationID = 1
	testFlavorID   = 1
)

func init() {
	resource.AddTestSweepers("serverscom_rbs_volume", &resource.Sweeper{
		Name: "serverscom_rbs_volume",
		F:    testSweepRBSVolumes,
	})
}

func testSweepRBSVolumes(region string) error {
	log.Printf("[DEBUG] Sweeping RBS volumes")
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping RBS volumes: %s", err)
	}

	ctx := context.TODO()
	volumes, err := client.RemoteBlockStorageVolumes.Collection().Collect(ctx)
	if err != nil {
		return fmt.Errorf("Error getting list of RBS volumes: %s", err)
	}

	for _, volume := range volumes {
		if !strings.HasPrefix(volume.Name, "tf-test-rbs-volume-") {
			continue
		}
		err := client.RemoteBlockStorageVolumes.Delete(ctx, volume.ID)
		if err != nil {
			return fmt.Errorf("Can't delete RBS volume (%s): %s", volume.ID, err)
		}
	}

	return nil
}

func TestAccServerscomRBSVolume_Basic(t *testing.T) {
	var rbsVolume scgo.RemoteBlockStorageVolume
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckRBSVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomRBSVolumeConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckRBSVolumeExists("serverscom_rbs_volume.test", &rbsVolume),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "name", fmt.Sprintf("tf-test-rbs-volume-%d", rInt)),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "size", "50"),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "status", "active"),
					resource.TestCheckResourceAttrSet(
						"serverscom_rbs_volume.test", "location_code"),
					resource.TestCheckResourceAttrSet(
						"serverscom_rbs_volume.test", "flavor_name"),
				),
			},
		},
	})
}

func TestAccServerscomRBSVolume_Update(t *testing.T) {
	var rbsVolume scgo.RemoteBlockStorageVolume
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckRBSVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomRBSVolumeConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckRBSVolumeExists("serverscom_rbs_volume.test", &rbsVolume),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "name", fmt.Sprintf("tf-test-rbs-volume-%d", rInt)),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "size", "10"),
				),
			},
			{
				Config: testAccServerscomRBSVolumeConfig_updated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckRBSVolumeExists("serverscom_rbs_volume.test", &rbsVolume),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "name", fmt.Sprintf("tf-test-rbs-volume-updated-%d", rInt)),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "size", "20"),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccServerscomRBSVolume_Resize(t *testing.T) {
	var rbsVolume scgo.RemoteBlockStorageVolume
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckRBSVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomRBSVolumeConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckRBSVolumeExists("serverscom_rbs_volume.test", &rbsVolume),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "size", "10"),
				),
			},
			{
				Config: testAccServerscomRBSVolumeConfig_resize(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckRBSVolumeExists("serverscom_rbs_volume.test", &rbsVolume),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "size", "30"),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccServerscomRBSVolume_WithLabels(t *testing.T) {
	var rbsVolume scgo.RemoteBlockStorageVolume
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckRBSVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomRBSVolumeConfig_withLabels(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckRBSVolumeExists("serverscom_rbs_volume.test", &rbsVolume),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "labels.environment", "test"),
					resource.TestCheckResourceAttr(
						"serverscom_rbs_volume.test", "labels.purpose", "acceptance-test"),
				),
			},
		},
	})
}

func testAccServerscomRBSVolumeConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "serverscom_rbs_volume" "test" {
	name        = "tf-test-rbs-volume-%d"
	size        = 50
	location_id = %d
	flavor_id   = %d
}
`, rInt, testLocationID, testFlavorID)
}

func testAccServerscomRBSVolumeConfig_updated(rInt int) string {
	return fmt.Sprintf(`
resource "serverscom_rbs_volume" "test" {
	name        = "tf-test-rbs-volume-updated-%d"
	size        = 50
	location_id = %d
	flavor_id   = %d
}
`, rInt, testLocationID, testFlavorID)
}

func testAccServerscomRBSVolumeConfig_resize(rInt int) string {
	return fmt.Sprintf(`
resource "serverscom_rbs_volume" "test" {
	name        = "tf-test-rbs-volume-%d"
	size        = 60
	location_id = %d
	flavor_id   = %d
}
`, rInt, testLocationID, testFlavorID)
}

func testAccServerscomRBSVolumeConfig_withLabels(rInt int) string {
	return fmt.Sprintf(`
resource "serverscom_rbs_volume" "test" {
	name        = "tf-test-rbs-volume-%d"
	size        = 50
	location_id = %d
	flavor_id   = %d
	
	labels = {
		environment = "test"
		purpose     = "acceptance-test"
	}
}
`, rInt, testLocationID, testFlavorID)
}

func testAccServerscomCheckRBSVolumeExists(n string, rbsVolume *scgo.RemoteBlockStorageVolume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RBS volume ID is set")
		}

		client := testAccProvider.Meta().(*scgo.Client)
		currentVolume, err := client.RemoteBlockStorageVolumes.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		*rbsVolume = *currentVolume
		return nil
	}
}

func testAccServerscomCheckRBSVolumeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*scgo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "serverscom_rbs_volume" {
			continue
		}

		_, err := client.RemoteBlockStorageVolumes.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			switch err.(type) {
			case *scgo.NotFoundError:
				return nil
			default:
				return fmt.Errorf("Error retrieving RBS volume: %s", err)
			}
		}

		return fmt.Errorf("RBS volume (%s) still exists", rs.Primary.ID)
	}

	return nil
}
