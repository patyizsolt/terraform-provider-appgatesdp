package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPolicyBasic(t *testing.T) {
	resourceName := "appgatesdp_policy.test_policy"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicyBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "expression", "return true;\n"),
					resource.TestCheckResourceAttr(resourceName, "notes", "terraform policy notes"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tamper_proofing", "true"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.0", "developer"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_links.0", "developer"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckPolicyBasic(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_policy" "test_policy" {
    name  = "%s"
    notes = "terraform policy notes"
    tags = [
        "terraform",
        "api-created"
    ]
    disabled = false

    expression = <<-EOF
        return true;
    EOF
    entitlement_links = [
        "developer"
    ]
    ringfence_rule_links = [
        "developer"
    ]
    tamper_proofing = true
    proxy_auto_config {
        enabled = true
        url     = "http://foo.com"
        persist = false
    }
    trusted_network_check {
        enabled    = true
        dns_suffix = "aa"
    }
}

`, rName)
}

func testAccCheckPolicyExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.PoliciesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.PoliciesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching policy with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgatesdp_policy" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.PoliciesApi

		if _, _, err := api.PoliciesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("policy still exists, %+v", err)
		}
	}
	return nil
}

//  TestAccPolicyClientSettings55 is only applicable on appliance 5.5
func TestAccPolicyClientSettings55(t *testing.T) {
	resourceName := "appgatesdp_policy.device_policy_with_client_settings"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name":         rName,
		"updated_name": "updated" + rName,
		"expression": `<<-EOF
		var result = false;
		return result;
		EOF`,
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c := testAccProvider.Meta().(*Client)
					c.GetToken()
					currentVersion := c.ApplianceVersion
					if currentVersion.LessThan(Appliance55Version) {
						t.Skip("Test only for 5.5 and above, appliance.portal is only supported in 5.4 and above.")
					}
				},
				Config: testAccCheckPolicyClientSettings(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.%", "10"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.add_remove_profiles", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.attention_level", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.auto_start", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.entitlements_list", "Hide"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.keep_me_signed_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.quit", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.saml_auto_sign_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.sign_out", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.suspend", "Show"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "expression", "var result = false;\nreturn result;\n"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "override_site_claim", ""),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.persist", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.url", ""),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tamper_proofing", "true"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.dns_suffix", ""),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "type", "Device"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckPolicyClientSettingsUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.%", "10"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.add_remove_profiles", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.attention_level", "High"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.auto_start", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.entitlements_list", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.keep_me_signed_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.quit", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.saml_auto_sign_in", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.sign_out", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.suspend", "Show"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["updated_name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "override_site_claim", ""),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.persist", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.url", ""),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "foo"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tamper_proofing", "true"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.dns_suffix", ""),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "type", "Device"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckPolicyClientSettings(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_policy" "device_policy_with_client_settings" {
	name = "%{name}"
	type = "Device"
	tags = [
		"terraform",
		"api-created"
	]
	disabled = false
	client_settings {
		enabled           = true
		entitlements_list = "Hide"
	}
	expression = <<-EOF
	var result = false;
	return result;
	EOF
}
`, context)
}

func testAccCheckPolicyClientSettingsUpdated(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_policy" "device_policy_with_client_settings" {
	name = "%{updated_name}"
	type = "Device"
	tags = [
		"terraform",
		"api-created",
		"foo"
	]
	disabled = false
	client_settings {
		enabled             = true
		entitlements_list   = "Show"
		attention_level     = "High"
		auto_start          = "Enabled"
		add_remove_profiles = "Show"
		keep_me_signed_in   = "Show"
		saml_auto_sign_in   = "Enabled"
		quit                = "Show"
		sign_out            = "Show"
		suspend             = "Show"
	}
	expression = <<-EOF
	var result = false;
	return result;
	EOF
}
`, context)
}
func TestAccPolicyDnsSettings55(t *testing.T) {
	resourceName := "appgatesdp_policy.dns_policy_with_dns_settings"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name":         rName,
		"updated_name": "updated" + rName,
		"expression": `<<-EOF
		var result = false;
		return result;
		EOF`,
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					c := testAccProvider.Meta().(*Client)
					c.GetToken()
					currentVersion := c.ApplianceVersion
					if currentVersion.LessThan(Appliance55Version) {
						t.Skip("Test only for 5.5 and above, appliance.portal is only supported in 5.4 and above.")
					}
				},
				Config: testAccCheckPolicyDnsSettings(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.%", "10"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.add_remove_profiles", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.attention_level", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.auto_start", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.entitlements_list", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.keep_me_signed_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.quit", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.saml_auto_sign_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.sign_out", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.suspend", "Show"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.domain", "appgate.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.1", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "override_site_claim", ""),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.persist", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.url", ""),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tamper_proofing", "true"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.dns_suffix", ""),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "type", "Dns"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckPolicyDnsSettingsUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "administrative_roles.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.%", "10"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.add_remove_profiles", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.attention_level", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.auto_start", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.entitlements_list", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.keep_me_signed_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.quit", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.saml_auto_sign_in", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.sign_out", "Show"),
					resource.TestCheckResourceAttr(resourceName, "client_settings.0.suspend", "Show"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.domain", "appgate.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.0.servers.1", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.1.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.1.domain", "google.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.1.servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.1.servers.0", "2.2.2.2"),
					resource.TestCheckResourceAttr(resourceName, "dns_settings.1.servers.1", "3.3.3.3"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlements.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["updated_name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "override_site_claim", ""),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.persist", "false"),
					resource.TestCheckResourceAttr(resourceName, "proxy_auto_config.0.url", ""),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rule_links.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ringfence_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "updated"),
					resource.TestCheckResourceAttr(resourceName, "tamper_proofing", "true"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.dns_suffix", ""),
					resource.TestCheckResourceAttr(resourceName, "trusted_network_check.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "type", "Dns"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckPolicyDnsSettings(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_policy" "dns_policy_with_dns_settings" {
	name = "%{name}"
	type = "Dns"
	tags = [
		"terraform",
		"api-created"
	]
	disabled = false
	dns_settings {
		domain  = "appgate.com"
		servers = ["8.8.8.8", "1.1.1.1"]
	}
	expression = <<-EOF
	var result = false;
	return result;
	EOF
}
`, context)
}

func testAccCheckPolicyDnsSettingsUpdated(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_policy" "dns_policy_with_dns_settings" {
	name = "%{updated_name}"
	type = "Dns"
	tags = [
		"terraform",
		"api-created",
		"updated",
	]
	disabled = false
	dns_settings {
		domain  = "appgate.com"
		servers = ["8.8.8.8", "1.1.1.1"]
	}
	dns_settings {
		domain  = "google.com"
		servers = ["2.2.2.2", "3.3.3.3"]
	}
	expression = <<-EOF
	var result = false;
	return result;
	EOF
}
`, context)
}
