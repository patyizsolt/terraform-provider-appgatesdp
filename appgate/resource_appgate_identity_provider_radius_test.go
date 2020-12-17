package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRadiusIdentityProviderBasic(t *testing.T) {
	resourceName := "appgate_radius_identity_provider.radius_test_resource"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRadiusIdentityProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRadiusIdentityProviderBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRadiusIdentityProviderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "admin_provider", "true"),
					resource.TestCheckResourceAttr(resourceName, "authentication_protocol", "CHAP"),
					resource.TestCheckResourceAttr(resourceName, "block_local_dns_requests", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2429955231.attribute_name", "memberOf"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2429955231.claim_name", "groups"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2429955231.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2429955231.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2440191359.attribute_name", "objectGUID"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2440191359.claim_name", "userId"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2440191359.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2440191359.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2672447578.attribute_name", "sn"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2672447578.claim_name", "lastName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2672447578.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.2672447578.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3609328139.attribute_name", "mail"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3609328139.claim_name", "emails"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3609328139.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3609328139.list", "true"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3889155743.attribute_name", "givenName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3889155743.claim_name", "firstName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3889155743.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.3889155743.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.633464825.attribute_name", "sAMAccountName"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.633464825.claim_name", "username"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.633464825.encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "claim_mappings.633464825.list", "false"),
					resource.TestCheckResourceAttr(resourceName, "default", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_search_domains.0", "internal.company.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.0", "172.17.18.19"),
					resource.TestCheckResourceAttr(resourceName, "dns_servers.1", "192.100.111.31"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "hostnames.0", "radius.company.com"),
					resource.TestCheckResourceAttr(resourceName, "inactivity_timeout_minutes", "0"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v4", "f572b4ab-7963-4a90-9e5a-3bf033bfe2cc"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_v6", "6935b379-205d-4fdd-847f-a0b5f14aff53"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.device_limit_per_user", "6"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.message", "welcome"),
					resource.TestCheckResourceAttr(resourceName, "on_boarding_two_factor.0.mfa_provider_id", "3ae98d53-c520-437f-99e4-451f936e6d2c"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.3657813038.claim_name", "antiVirusIsRunning"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.3657813038.command", "fileSize"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.3657813038.parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.3657813038.parameters.0.args", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.3657813038.parameters.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.3657813038.parameters.0.path", "/usr/bin/python3"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_claim_mappings.3657813038.platform", "desktop.windows.all"),
					resource.TestCheckResourceAttr(resourceName, "port", "1812"),
					resource.TestCheckResourceAttr(resourceName, "shared_secret", "hunter2"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "type", "Radius"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccRadiusIdentityProviderImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"shared_secret"},
			},
		},
	})
}

func testAccCheckRadiusIdentityProviderBasic(rName string) string {
	return fmt.Sprintf(`
data "appgate_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgate_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}
data "appgate_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgate_radius_identity_provider" "radius_test_resource" {
  name = "%s"
  hostnames = [
    "radius.company.com"
  ]
  admin_provider = true
  port           = 1812
  shared_secret  = "hunter2"
  ip_pool_v4     = data.appgate_ip_pool.ip_v4_pool.id
  ip_pool_v6     = data.appgate_ip_pool.ip_v6_pool.id
  dns_servers = [
    "172.17.18.19",
    "192.100.111.31"
  ]
  dns_search_domains = [
    "internal.company.com"
  ]
  block_local_dns_requests = true
  on_boarding_two_factor {
    mfa_provider_id       = data.appgate_mfa_provider.fido.id
    device_limit_per_user = 6
    message               = "welcome"
  }
  tags = [
    "terraform",
    "api-created"
  ]
  claim_mappings {
    attribute_name = "objectGUID"
    claim_name     = "userId"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sAMAccountName"
    claim_name     = "username"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "givenName"
    claim_name     = "firstName"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sn"
    claim_name     = "lastName"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "mail"
    claim_name     = "emails"
    encrypted      = false
    list           = true
  }
  claim_mappings {
    attribute_name = "memberOf"
    claim_name     = "groups"
    encrypted      = false
    list           = true
  }

  on_demand_claim_mappings {
    command    = "fileSize"
    claim_name = "antiVirusIsRunning"
    parameters {
      path = "/usr/bin/python3"
    }
    platform = "desktop.windows.all"
  }
}
`, rName)
}

func testAccCheckRadiusIdentityProviderExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.RadiusIdentityProvidersApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching radius identity provider with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckRadiusIdentityProviderDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_radius_identity_provider" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.RadiusIdentityProvidersApi

		_, _, err := api.IdentityProvidersIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("radius identity provider still exists, %+v", err)
		}
	}
	return nil
}

func testAccRadiusIdentityProviderImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
