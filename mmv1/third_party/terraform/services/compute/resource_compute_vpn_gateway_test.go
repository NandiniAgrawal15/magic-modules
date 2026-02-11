package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeVpnGateway_resourceManagerTags(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	suffix := acctest.RandString(t, 10)
	tagKeyResult := acctest.BootstrapSharedTestTagKeyDetails(t, "crm-vpn-tagkey", "organizations/"+org, make(map[string]interface{}))
	sharedTagkey, _ := tagKeyResult["shared_tag_key"]
	tagValueResult := acctest.BootstrapSharedTestTagValueDetails(t, "crm-vpn-tagvalue", sharedTagkey, org)

	context := map[string]interface{}{
		"suffix":       suffix,
		"tag_key_id":   tagKeyResult["name"],
		"tag_value_id": tagValueResult["name"],
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVpnGateway_resourceManagerTags(context),
			},
			{
				ResourceName:            "google_compute_vpn_gateway.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"params"},
			},
		},
	})
}

func testAccComputeVpnGateway_resourceManagerTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-network-%{suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "tf-test-vpn-%{suffix}"
  network = google_compute_network.foobar.id
  region  = "us-central1"
  params {
    resource_manager_tags = {
      "%{tag_key_id}" = "%{tag_value_id}"
    }
  }
}
`, context)
}
