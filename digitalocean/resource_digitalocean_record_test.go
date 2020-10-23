package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestDigitalOceanRecordConstructFqdn(t *testing.T) {
	cases := []struct {
		Input, Output string
	}{
		{"www", "www.nonexample.com"},
		{"dev.www", "dev.www.nonexample.com"},
		{"*", "*.nonexample.com"},
		{"nonexample.com", "nonexample.com.nonexample.com"},
		{"test.nonexample.com", "test.nonexample.com.nonexample.com"},
		{"test.nonexample.com.", "test.nonexample.com"},
	}

	domain := "nonexample.com"
	for _, tc := range cases {
		actual := constructFqdn(tc.Input, domain)
		if actual != tc.Output {
			t.Fatalf("input: %s\noutput: %s", tc.Input, actual)
		}
	}
}

func TestAccDigitalOceanRecord_Basic(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "fqdn", strings.Join([]string{"terraform", domain}, ".")),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_BasicFullName(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanRecordConfig_basic_full_name, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "fqdn", strings.Join([]string{"terraform", domain}, ".")),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_Updated(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "type", "A"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "ttl", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_new_value, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "type", "A"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "ttl", "90"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_HostnameValue(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_cname, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("a.foobar-test-terraform.com", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "a.foobar-test-terraform.com."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "type", "CNAME"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_ExternalHostnameValue(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_external_cname, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("a.foobar-test-terraform.net", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "a.foobar-test-terraform.net."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "type", "CNAME"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_FlagsAndTag(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_caa, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("letsencrypt.org", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "letsencrypt.org."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "type", "CAA"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "flags", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "tag", "issue"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_MX(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_mx, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "MX"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_MX_at(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_mx_at, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("@", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "@"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "MX"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_SRV_zero_weight_port(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_srv_zero_weight_port, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "_service._protocol"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "SRV"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "weight", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "port", "0"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_UpdateBasic(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_updated_basic, domain, "terraform", "a.foobar-test-terraform.com.", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("a.foobar-test-terraform.com", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "a.foobar-test-terraform.com."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "type", "CNAME"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "ttl", "1800"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_updated_basic, domain, "terraform-updated", "b.foobar-test-terraform.com.", "1000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foobar", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("b.foobar-test-terraform.com", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "name", "terraform-updated"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "value", "b.foobar-test-terraform.com."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "type", "CNAME"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foobar", "ttl", "1000"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_MXUpdated(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_mx_updated, domain, "10"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "MX"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "priority", "10"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_mx_updated, domain, "20"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "MX"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "priority", "20"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_SrvUpdated(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_srv_updated, domain, "5050", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "_service._protocol"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "SRV"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "port", "5050"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "weight", "100"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_srv_updated, domain, "6060", "150"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("foobar."+domain, &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "_service._protocol"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "foobar."+domain+"."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "SRV"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "port", "6060"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "weight", "150"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_CaaUpdated(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_caa_updated, domain, "20", "issue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("letsencrypt.org", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "letsencrypt.org."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "CAA"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "flags", "20"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "tag", "issue"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_caa_updated, domain, "50", "issuewild"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.foo_record", &record),
					testAccCheckDigitalOceanRecordAttributesHostname("letsencrypt.org", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "value", "letsencrypt.org."),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "type", "CAA"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "flags", "50"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.foo_record", "tag", "issuewild"),
				),
			},
		},
	})
}

func TestAccDigitalOceanRecord_iodefCAA(t *testing.T) {
	var record godo.DomainRecord
	domain := fmt.Sprintf("foobar-test-terraform-%s.com", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckDigitalOceanRecordConfig_iodef, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanRecordExists("digitalocean_record.CAA_iodef", &record),
					resource.TestCheckResourceAttr(
						"digitalocean_record.CAA_iodef", "name", "@"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.CAA_iodef", "domain", domain),
					resource.TestCheckResourceAttr(
						"digitalocean_record.CAA_iodef", "value", "mailto:caa-failures@example.com"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.CAA_iodef", "type", "CAA"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.CAA_iodef", "flags", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_record.CAA_iodef", "tag", "iodef"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_record" {
			continue
		}
		domain := rs.Primary.Attributes["domain"]
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, _, err = client.Domains.Record(context.Background(), domain, id)

		if err == nil {
			return fmt.Errorf("Record still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanRecordAttributes(record *godo.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Data != "192.168.0.10" {
			return fmt.Errorf("Bad value: %s", record.Data)
		}

		return nil
	}
}

func testAccCheckDigitalOceanRecordAttributesUpdated(record *godo.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Data != "192.168.0.11" {
			return fmt.Errorf("Bad value: %s", record.Data)
		}

		return nil
	}
}

func testAccCheckDigitalOceanRecordExists(n string, record *godo.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		domain := rs.Primary.Attributes["domain"]
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		foundRecord, _, err := client.Domains.Record(context.Background(), domain, id)

		if err != nil {
			return err
		}

		if strconv.Itoa(foundRecord.ID) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*record = *foundRecord

		return nil
	}
}

func testAccCheckDigitalOceanRecordAttributesHostname(data string, record *godo.DomainRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Data != data {
			return fmt.Errorf("Bad value: expected %s, got %s", data, record.Data)
		}

		return nil
	}
}

const testAccCheckDigitalOceanRecordConfig_basic = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foobar" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "terraform"
  value = "192.168.0.10"
  type  = "A"
}`

const testAccCheckDigitalOceanRecordConfig_basic_full_name = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foobar" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "terraform.${digitalocean_domain.foobar.name}."
  value = "192.168.0.10"
  type  = "A"
}`

const testAccCheckDigitalOceanRecordConfig_new_value = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foobar" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "terraform"
  value = "192.168.0.11"
  type  = "A"
  ttl   = 90
}`

const testAccCheckDigitalOceanRecordConfig_cname = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foobar" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "terraform"
  value = "a.foobar-test-terraform.com."
  type  = "CNAME"
}`

const testAccCheckDigitalOceanRecordConfig_mx_at = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foo_record" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "terraform"
  value = "${digitalocean_domain.foobar.name}."
  type  = "MX"
  priority = "10"
}`

const testAccCheckDigitalOceanRecordConfig_mx = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foo_record" {
  domain = "${digitalocean_domain.foobar.name}"

  name     = "terraform"
  value    = "foobar.${digitalocean_domain.foobar.name}."
  type     = "MX"
  priority = "10"
}`

const testAccCheckDigitalOceanRecordConfig_external_cname = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foobar" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "terraform"
  value = "a.foobar-test-terraform.net."
  type  = "CNAME"
}`

const testAccCheckDigitalOceanRecordConfig_caa = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foobar" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "terraform"
  type  = "CAA"
  value = "letsencrypt.org."
  flags = 1
  tag   = "issue"
}`

const testAccCheckDigitalOceanRecordConfig_srv_zero_weight_port = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foo_record" {
  domain = "${digitalocean_domain.foobar.name}"

  name     = "_service._protocol"
  value    = "foobar.${digitalocean_domain.foobar.name}."
  type     = "SRV"
  priority = 10
  port     = 0
  weight   = 0
}`

const testAccCheckDigitalOceanRecordConfig_updated_basic = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foobar" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "%s"
  value = "%s"
  type  = "CNAME"
  ttl   = "%s"
}`

const testAccCheckDigitalOceanRecordConfig_mx_updated = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foo_record" {
  domain = "${digitalocean_domain.foobar.name}"

  name     = "terraform"
  value    = "foobar.${digitalocean_domain.foobar.name}."
  type     = "MX"
  priority = "%s"
}`

const testAccCheckDigitalOceanRecordConfig_srv_updated = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foo_record" {
  domain = "${digitalocean_domain.foobar.name}"

  name     = "_service._protocol"
  value    = "foobar.${digitalocean_domain.foobar.name}."
  type     = "SRV"
  priority = "10"
  port     = "%s"
  weight   = "%s"
}`

const testAccCheckDigitalOceanRecordConfig_caa_updated = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
  ip_address = "192.168.0.10"
}

resource "digitalocean_record" "foo_record" {
  domain = "${digitalocean_domain.foobar.name}"

  name  = "terraform"
  type  = "CAA"
  value = "letsencrypt.org."
  flags = "%s"
  tag   = "%s"
}`

const testAccCheckDigitalOceanRecordConfig_iodef = `
resource "digitalocean_domain" "foobar" {
  name       = "%s"
}
resource "digitalocean_record" "CAA_iodef" {
  domain = digitalocean_domain.foobar.name
  type  = "CAA"
  tag   = "iodef"
  flags = "0"
  name  = "@"
  value = "mailto:caa-failures@example.com"
}`
