package serverscom

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

// testAccDNSDomainBase returns the parent domain the acceptance tests build
// their zones under. Override it with SERVERSCOM_TEST_DNS_DOMAIN to run against
// a domain you control (e.g. so delegation actually resolves); it defaults to
// example.com. Generated zone names keep the tf-test- prefix so the sweeper
// still reclaims them.
func testAccDNSDomainBase() string {
	if v := os.Getenv("SERVERSCOM_TEST_DNS_DOMAIN"); v != "" {
		return v
	}
	return "example.com"
}

func init() {
	resource.AddTestSweepers("serverscom_dns_domain", &resource.Sweeper{
		Name: "serverscom_dns_domain",
		F:    testSweepDNSDomains,
	})
}

func testSweepDNSDomains(region string) error {
	log.Printf("[DEBUG] Sweeping DNS domains")
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("error getting client for sweeping DNS domains: %s", err)
	}

	ctx := context.TODO()
	domains, err := client.DNS.Collection().Collect(ctx)
	if err != nil {
		return fmt.Errorf("error getting list of DNS domains: %s", err)
	}

	for _, domain := range domains {
		if !strings.HasPrefix(domain.Name, "tf-test-") {
			continue
		}
		if err := client.DNS.DeleteDomain(ctx, domain.ID); err != nil {
			return fmt.Errorf("can't delete DNS domain (%s): %s", domain.ID, err)
		}
	}

	return nil
}

func TestAccServerscomDNSDomain_Basic(t *testing.T) {
	var domain scgo.DNSDomain
	rInt := acctest.RandInt()
	name := fmt.Sprintf("tf-test-%d.%s", rInt, testAccDNSDomainBase())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomDNSDomainConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDNSDomainExists("serverscom_dns_domain.test", &domain),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain.test", "name", name),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain.test", "email", "admin@example.com"),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain.test", "ttl", "3600"),
					resource.TestCheckResourceAttrSet(
						"serverscom_dns_domain.test", "delegation_status"),
					resource.TestCheckResourceAttrSet(
						"serverscom_dns_domain.test", "created_at"),
				),
			},
			{
				ResourceName:      "serverscom_dns_domain.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServerscomDNSDomain_WithLabels(t *testing.T) {
	var domain scgo.DNSDomain
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomDNSDomainConfig_withLabels(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDNSDomainExists("serverscom_dns_domain.test", &domain),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain.test", "labels.env", "test"),
				),
			},
			{
				Config: testAccServerscomDNSDomainConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDNSDomainExists("serverscom_dns_domain.test", &domain),
					resource.TestCheckNoResourceAttr(
						"serverscom_dns_domain.test", "labels.env"),
				),
			},
		},
	})
}

func testAccServerscomDNSDomainConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "serverscom_dns_domain" "test" {
	name  = "tf-test-%d.%s"
	email = "admin@example.com"
	ttl   = 3600
}
`, rInt, testAccDNSDomainBase())
}

func testAccServerscomDNSDomainConfig_withLabels(rInt int) string {
	return fmt.Sprintf(`
resource "serverscom_dns_domain" "test" {
	name  = "tf-test-%d.%s"
	email = "admin@example.com"
	ttl   = 3600

	labels = {
		env = "test"
	}
}
`, rInt, testAccDNSDomainBase())
}

func testAccServerscomCheckDNSDomainExists(n string, domain *scgo.DNSDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no DNS domain ID is set")
		}

		client := testAccProvider.Meta().(*scgo.Client)
		current, err := client.DNS.GetDomain(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		*domain = *current
		return nil
	}
}

func testAccServerscomCheckDNSDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*scgo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "serverscom_dns_domain" {
			continue
		}

		_, err := client.DNS.GetDomain(context.Background(), rs.Primary.ID)
		if err != nil {
			switch err.(type) {
			case *scgo.NotFoundError:
				continue
			default:
				return fmt.Errorf("error retrieving DNS domain: %s", err)
			}
		}

		return fmt.Errorf("DNS domain (%s) still exists", rs.Primary.ID)
	}

	return nil
}
