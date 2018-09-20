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

func TestAccMonitoringGroup_basic(t *testing.T) {
	t.Parallel()

	groupName := acctest.RandomWithPrefix("tf-test-monitoring-group")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMonitoringGroup_basic(groupName),
			},
			resource.TestStep{
				ResourceName:            "google_monitoring_group.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func testAccCheckMonitoringGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_monitoring_group" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		group, _ := config.clientMonitoringGroup.GetGroup(context.Background(), &monitoring.GetGroupRequest{Name: rs.Primary.ID})
		if group != nil {
			return fmt.Errorf("Group still present")
		}
	}

	return nil
}

func testAccMonitoringGroup_basic(displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_group" "foo" {
	display_name = "%s"
  filter = "resource.metadata.region=\"europe-west2\""
}`, displayName)
}
