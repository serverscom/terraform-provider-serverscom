package serverscom

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

func TestAccServerscomDNSDomainRecord_Basic(t *testing.T) {
	var record scgo.DNSRecord
	rInt := acctest.RandInt()
	name := fmt.Sprintf("www.tf-test-%d.%s", rInt, testAccDNSDomainBase())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomDNSDomainRecordConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDNSDomainRecordExists("serverscom_dns_domain_record.a", &record),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.a", "name", name),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.a", "type", "A"),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.a", "data", "192.0.2.1"),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.a", "ttl", "300"),
					resource.TestCheckResourceAttrSet(
						"serverscom_dns_domain_record.a", "created_at"),
				),
			},
			{
				ResourceName:      "serverscom_dns_domain_record.a",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccServerscomDNSDomainRecordImportID("serverscom_dns_domain_record.a"),
			},
		},
	})
}

func TestAccServerscomDNSDomainRecord_Update(t *testing.T) {
	var record scgo.DNSRecord
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomDNSDomainRecordConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDNSDomainRecordExists("serverscom_dns_domain_record.a", &record),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.a", "data", "192.0.2.1"),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.a", "ttl", "300"),
				),
			},
			{
				Config: testAccServerscomDNSDomainRecordConfig_updated(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDNSDomainRecordExists("serverscom_dns_domain_record.a", &record),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.a", "data", "192.0.2.2"),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.a", "ttl", "600"),
				),
			},
		},
	})
}

func TestAccServerscomDNSDomainRecord_MX(t *testing.T) {
	var record scgo.DNSRecord
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccServerscomPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServerscomCheckDNSDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerscomDNSDomainRecordConfig_mx(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccServerscomCheckDNSDomainRecordExists("serverscom_dns_domain_record.mx", &record),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.mx", "type", "MX"),
					resource.TestCheckResourceAttr(
						"serverscom_dns_domain_record.mx", "priority", "10"),
				),
			},
		},
	})
}

func testAccServerscomDNSDomainRecordConfig_basic(rInt int) string {
	base := testAccDNSDomainBase()
	return fmt.Sprintf(`
resource "serverscom_dns_domain" "test" {
	name  = "tf-test-%d.%s"
	email = "admin@example.com"
	ttl   = 3600
}

resource "serverscom_dns_domain_record" "a" {
	dns_domain_id = serverscom_dns_domain.test.id
	name          = "www.tf-test-%d.%s"
	type          = "A"
	data          = "192.0.2.1"
	ttl           = 300
}
`, rInt, base, rInt, base)
}

func testAccServerscomDNSDomainRecordConfig_updated(rInt int) string {
	base := testAccDNSDomainBase()
	return fmt.Sprintf(`
resource "serverscom_dns_domain" "test" {
	name  = "tf-test-%d.%s"
	email = "admin@example.com"
	ttl   = 3600
}

resource "serverscom_dns_domain_record" "a" {
	dns_domain_id = serverscom_dns_domain.test.id
	name          = "www.tf-test-%d.%s"
	type          = "A"
	data          = "192.0.2.2"
	ttl           = 600
}
`, rInt, base, rInt, base)
}

func testAccServerscomDNSDomainRecordConfig_mx(rInt int) string {
	base := testAccDNSDomainBase()
	return fmt.Sprintf(`
resource "serverscom_dns_domain" "test" {
	name  = "tf-test-%d.%s"
	email = "admin@example.com"
	ttl   = 3600
}

resource "serverscom_dns_domain_record" "mx" {
	dns_domain_id = serverscom_dns_domain.test.id
	name          = "tf-test-%d.%s"
	type          = "MX"
	data          = "mail.tf-test-%d.%s"
	priority      = 10
	ttl           = 3600
}
`, rInt, base, rInt, base, rInt, base)
}

func testAccServerscomDNSDomainRecordImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["dns_domain_id"], rs.Primary.ID), nil
	}
}

func testAccServerscomCheckDNSDomainRecordExists(n string, record *scgo.DNSRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no DNS record ID is set")
		}

		domainID := rs.Primary.Attributes["dns_domain_id"]

		client := testAccProvider.Meta().(*scgo.Client)
		current, err := client.DNS.GetRecord(context.Background(), domainID, rs.Primary.ID)
		if err != nil {
			return err
		}

		*record = *current
		return nil
	}
}
