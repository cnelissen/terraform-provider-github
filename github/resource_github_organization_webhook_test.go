package github

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccGithubOrganizationWebhook(t *testing.T) {

	randomID := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	t.Run("creates and updates webhooks without error", func(t *testing.T) {

		config := fmt.Sprintf(`

			resource "github_repository" "test" {
			  name = "tf-acc-test-%s"
				auto_init = true
			}

			resource "github_organization_webhook" "test" {
			  configuration {
			    url = "https://google.de/webhook"
			    content_type = "json"
			    insecure_ssl = true
			  }

			  events = ["pull_request"]
			}

		`, randomID)

		checks := map[string]resource.TestCheckFunc{
			"before": resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(
					"github_organization_webhook.test", "events.0",
					"pull_request",
				),
			),
			"after": resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(
					"github_organization_webhook.test", "events.0",
					"issues",
				),
			),
		}

		testCase := func(t *testing.T, mode string) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { skipUnlessMode(t, mode) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: config,
						Check:  checks["before"],
					},
					{
						Config: strings.Replace(config,
							"https://google.de/webhook",
							"https://google.de/updated", 1),
						Check: checks["after"],
					},
				},
			})
		}

		t.Run("with an anonymous account", func(t *testing.T) {
			t.Skip("anonymous account not supported for this operation")
		})

		t.Run("with an individual account", func(t *testing.T) {
			t.Skip("individual account not supported for this operation")
		})

		t.Run("with an organization account", func(t *testing.T) {
			testCase(t, organization)
		})
	})

	// t.Run("imports webhooks without error", func(t *testing.T) {
	//
	// 	config := fmt.Sprintf(`
	//
	// 		resource "github_repository" "test" {
	// 		  name = "tf-acc-test-%s"
	// 			auto_init = true
	// 		}
	//
	// 		resource "github_organization_webhook" "test" {
	// 		  configuration {
	// 		    url = "https://google.de/webhook"
	// 		    content_type = "json"
	// 		    insecure_ssl = true
	// 		  }
	//
	// 		  events = ["issue"]
	// 		}
	//
	// 	`, randomID)
	//
	// 	check := resource.ComposeTestCheckFunc(
	// 		resource.TestCheckResourceAttr(
	// 			"github_organization_webhook.test", "events.0",
	// 			"issue",
	// 		),
	// 	)
	//
	// 	testCase := func(t *testing.T, mode string) {
	// 		resource.Test(t, resource.TestCase{
	// 			PreCheck:  func() { skipUnlessMode(t, mode) },
	// 			Providers: testAccProviders,
	// 			Steps: []resource.TestStep{
	// 				{
	// 					Config: config,
	// 					Check:  check,
	// 				},
	// 				{
	// 					ResourceName:      "github_organization_webhook.test",
	// 					ImportState:       true,
	// 					ImportStateVerify: true,
	// 					// ImportStateIdPrefix: fmt.Sprintf("foo-%s/", randString),
	// 				},
	// 			},
	// 		})
	// 	}
	//
	// 	t.Run("with an anonymous account", func(t *testing.T) {
	// 		t.Skip("anonymous account not supported for this operation")
	// 	})
	//
	// 	t.Run("with an individual account", func(t *testing.T) {
	// 		t.Skip("individual account not supported for this operation")
	// 	})
	//
	// 	t.Run("with an organization account", func(t *testing.T) {
	// 		testCase(t, organization)
	// 	})
	// })

}

// func TestAccGithubOrganizationWebhook_basic(t *testing.T) {
// 	if err := testAccCheckOrganization(); err != nil {
// 		t.Skipf("Skipping because %s.", err.Error())
// 	}
//
// 	var hook github.Hook
// 	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
// 	rn := "github_organization_webhook.foo"
//
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckGithubOrganizationWebhookDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccGithubOrganizationWebhookConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckGithubOrganizationWebhookExists(rn, &hook),
// 					testAccCheckGithubOrganizationWebhookAttributes(&hook, &testAccGithubOrganizationWebhookExpectedAttributes{
// 						Events: []string{"pull_request"},
// 						Configuration: map[string]interface{}{
// 							"url":          "https://google.de/webhook",
// 							"content_type": "json",
// 							"insecure_ssl": "1",
// 						},
// 						Active: true,
// 					}),
// 				),
// 			},
// 			{
// 				Config: testAccGithubOrganizationWebhookUpdateConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckGithubOrganizationWebhookExists(rn, &hook),
// 					testAccCheckGithubOrganizationWebhookAttributes(&hook, &testAccGithubOrganizationWebhookExpectedAttributes{
// 						Events: []string{"issues"},
// 						Configuration: map[string]interface{}{
// 							"url":          "https://google.de/webhooks",
// 							"content_type": "form",
// 							"insecure_ssl": "0",
// 						},
// 						Active: false,
// 					}),
// 				),
// 			},
// 			{
// 				ResourceName:        rn,
// 				ImportState:         true,
// 				ImportStateVerify:   true,
// 				ImportStateIdPrefix: fmt.Sprintf("foo-%s/", randString),
// 			},
// 		},
// 	})
// }
//
// func TestAccGithubOrganizationWebhook_secret(t *testing.T) {
// 	if err := testAccCheckOrganization(); err != nil {
// 		t.Skipf("Skipping because %s.", err.Error())
// 	}
//
// 	rn := "github_organization_webhook.foo"
// 	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
//
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckGithubOrganizationWebhookDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccGithubOrganizationWebhookConfig_secret,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckGithubOrganizationWebhookSecret(rn, "VerySecret"),
// 				),
// 			},
// 			{
// 				ResourceName:        rn,
// 				ImportState:         true,
// 				ImportStateVerify:   true,
// 				ImportStateIdPrefix: fmt.Sprintf("foo-%s/", randString),
// 			},
// 		},
// 	})
// }
//
// func testAccCheckGithubOrganizationWebhookExists(n string, hook *github.Hook) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[n]
// 		if !ok {
// 			return fmt.Errorf("Not Found: %s", n)
// 		}
//
// 		hookID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
// 		if err != nil {
// 			return unconvertibleIdErr(rs.Primary.ID, err)
// 		}
// 		if hookID == 0 {
// 			return fmt.Errorf("No repository name is set")
// 		}
//
// 		org := testAccProvider.Meta().(*Owner)
// 		conn := org.v3client
// 		getHook, _, err := conn.Organizations.GetHook(context.TODO(), org.name, hookID)
// 		if err != nil {
// 			return err
// 		}
// 		*hook = *getHook
// 		return nil
// 	}
// }
//
// type testAccGithubOrganizationWebhookExpectedAttributes struct {
// 	Events        []string
// 	Configuration map[string]interface{}
// 	Active        bool
// }
//
// func testAccCheckGithubOrganizationWebhookAttributes(hook *github.Hook, want *testAccGithubOrganizationWebhookExpectedAttributes) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
//
// 		if active := hook.GetActive(); active != want.Active {
// 			return fmt.Errorf("got hook %t; want %t", active, want.Active)
// 		}
// 		if URL := hook.GetURL(); !strings.HasPrefix(URL, "https://") {
// 			return fmt.Errorf("got http URL %q; want to start with 'https://'", URL)
// 		}
// 		if !reflect.DeepEqual(hook.Events, want.Events) {
// 			return fmt.Errorf("got hook events %q; want %q", hook.Events, want.Events)
// 		}
// 		if !reflect.DeepEqual(hook.Config, want.Configuration) {
// 			return fmt.Errorf("got hook configuration %q; want %q", hook.Config, want.Configuration)
// 		}
//
// 		return nil
// 	}
// }
//
// func testAccCheckGithubOrganizationWebhookSecret(r, secret string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[r]
// 		if !ok {
// 			return fmt.Errorf("Not Found: %s", r)
// 		}
//
// 		if rs.Primary.Attributes["configuration.0.secret"] != secret {
// 			return fmt.Errorf("Configured secret in %s does not match secret in state.  (Expected: %s, Actual: %s)", r, secret, rs.Primary.Attributes["configuration.0.secret"])
// 		}
//
// 		return nil
// 	}
// }
//
// func testAccCheckGithubOrganizationWebhookDestroy(s *terraform.State) error {
// 	conn := testAccProvider.Meta().(*Owner).v3client
// 	orgName := testAccProvider.Meta().(*Owner).name
//
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "github_organization_webhook" {
// 			continue
// 		}
//
// 		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
// 		if err != nil {
// 			return unconvertibleIdErr(rs.Primary.ID, err)
// 		}
//
// 		gotHook, resp, err := conn.Organizations.GetHook(context.TODO(), orgName, id)
// 		if err == nil {
// 			if gotHook != nil && gotHook.GetID() == id {
// 				return fmt.Errorf("Webhook still exists")
// 			}
// 		}
// 		if resp.StatusCode != 404 {
// 			return err
// 		}
// 		return nil
// 	}
// 	return nil
// }
//
// const testAccGithubOrganizationWebhookConfig = `
// resource "github_organization_webhook" "foo" {
//   configuration {
//     url = "https://google.de/webhook"
//     content_type = "json"
//     insecure_ssl = true
//   }
//
//   events = ["pull_request"]
// }
// `
//
// const testAccGithubOrganizationWebhookUpdateConfig = `
// resource "github_organization_webhook" "foo" {
//   configuration {
//     url = "https://google.de/webhooks"
//     content_type = "form"
//     insecure_ssl = false
//   }
//   active = false
//
//   events = ["issues"]
// }
// `
//
// const testAccGithubOrganizationWebhookConfig_secret = `
// resource "github_organization_webhook" "foo" {
//   configuration {
//     url          = "https://www.terraform.io/webhook"
//     content_type = "json"
//     secret       = "VerySecret"
//     insecure_ssl = false
//   }
//
//   events = ["pull_request"]
// }
// `
