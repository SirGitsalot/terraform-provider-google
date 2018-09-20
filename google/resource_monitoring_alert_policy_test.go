package google

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/genproto/googleapis/monitoring/v3"
	"testing"
)

func TestAccMonitoringAlertPolicy_basic(t *testing.T) {
	t.Parallel()

	alertPolicyName := acctest.RandomWithPrefix("tf-test-monitoring-alert-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringAlertPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMonitoringAlertPolicy_basic(alertPolicyName),
			},
			resource.TestStep{
				ResourceName:            "google_monitoring_alert_policy.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func testAccCheckMonitoringAlertPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_monitoring_alert_policy" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		alertPolicy, _ := config.clientMonitoringPolicy.GetAlertPolicy(context.Background(), &monitoring.GetAlertPolicyRequest{Name: rs.Primary.ID})
		if alertPolicy != nil {
			return fmt.Errorf("AlertPolicy still present")
		}
	}

	return nil
}

func testAccMonitoringAlertPolicy_basic(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_alert_policy" "foo" {
  display_name = "%s"
  combiner = "OR"
  enabled = true
  conditions = [
    {
      display_name = "Test Condition"
      condition_threshold = {
        comparison = "COMPARISON_GT"
        threshold_value = 100000
        trigger = {
          count = 1
        }
        duration = "60s"
        filter = "resource.type=\"gcs_bucket\" AND metric.type=\"storage.googleapis.com/storage/total_bytes\""
      }
   }
  ]
}`, displayName)
}
