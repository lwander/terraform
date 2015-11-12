package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/logging/v1beta3"
)

func TestAccLoggingSink_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingSinkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccLoggingSink_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingSinkExists(
						"google_compute_logs_sink.foobar"),
				),
			},
		},
	})
}

func testAccCheckLoggingSinkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	logsSinksService := logging.NewProjectsLogsSinksService(config.clientLogging)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_logs_sink" {
			continue
		}

		sink_name := rs.Primary.Attributes["sink_name"]
		log_name := rs.Primary.Attributes["log_name"]

		_, err := logsSinksService.Get(config.Project, log_name, sink_name).Do()
		if err == nil {
			return fmt.Errorf("Sink %s still exists", sink_name)
		}
	}

	return nil
}

func testAccCheckLoggingSinkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		logsSinksService := logging.NewProjectsLogsSinksService(config.clientLogging)

		sink_name := rs.Primary.Attributes["sink_name"]
		log_name := rs.Primary.Attributes["log_name"]

		_, err := logsSinksService.Get(config.Project, log_name, sink_name).Do()
		if err != nil {
			return fmt.Errorf("Sink %s does not exist", sink_name)
		}

		return nil
	}
}

const testAccLoggingSink_basic = `
resource "google_storage_bucket" "foobar" {
	name = "tf-test-storage-bucket-v2"
}

resource "google_storage_bucket_acl" "everyone" {
	bucket = "${google_storage_bucket.foobar.name}"
	role_entity = ["OWNER:allUsers"]
}

resource "google_logging_sink" "foobar" {
	log_name = "activity_log"
	sink_name = "tf-test-basic"
	destination = "storage.googleapis.com/${google_storage_bucket.foobar.name}"
	depends_on = ["google_storage_bucket_acl.everyone"]
}`
