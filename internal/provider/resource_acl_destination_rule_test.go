// Copyright 2021 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tsuru/terraform-provider-acl/internal/acl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/tsuru/acl-api/api/types"
)

var testAccProvider *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"acl": func() (*schema.Provider, error) {
		return testAccProvider, nil
	},
}

func init() {
	testAccProvider = Provider()
}

func TestAccResourceDestinationRuleApp(t *testing.T) {
	fakeServer := echo.New()
	myRule := types.ServiceRule{
		Rule: types.Rule{
			RuleID: "my-rule",
			Destination: types.RuleType{
				TsuruApp: &types.TsuruAppRule{
					AppName: "my-destination-app",
				},
			},
		},
	}

	fakeServer.Any("/services/acl/proxy/:instance", func(c echo.Context) error {
		callback := c.QueryParam("callback")
		if callback == "/rule/"+myRule.RuleID && c.Request().Method == http.MethodDelete {
			return c.String(http.StatusOK, "")
		}

		if callback == "/rule" {
			if c.Request().Method == http.MethodPost {
				return c.JSON(http.StatusOK, myRule)
			}

			return c.JSON(http.StatusOK, &acl.ServiceRuleData{
				ServiceInstance: types.ServiceInstance{
					BaseRules: []types.ServiceRule{
						myRule,
					},
				},
			})
		}
		t.Fatalf("method=%q, path=%q, callback=%q, err=\"Not found\"",
			c.Request().Method,
			c.Path(),
			callback,
		)
		return c.String(http.StatusNotFound, "")
	})

	fakeServer.HTTPErrorHandler = func(err error, c echo.Context) {
		t.Errorf("methods=%s, path=%s, err=%s", c.Request().Method, c.Path(), err.Error())
	}
	server := httptest.NewServer(fakeServer)
	os.Setenv("TSURU_TARGET", server.URL)

	resourceName := "acl_destination_rule.rule"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: `
resource "acl_destination_rule" "rule" {
	instance =  "my-acl"

	app = "my-destination-app"
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app", "my-destination-app")),
			},
		},
	})
}

func TestAccImportRuleApp(t *testing.T) {
	fakeServer := echo.New()
	myRule := types.ServiceRule{
		Rule: types.Rule{
			RuleID: "my-rule",
			Destination: types.RuleType{
				TsuruApp: &types.TsuruAppRule{
					AppName: "my-destination-app",
				},
			},
		},
	}

	fakeServer.Any("/services/acl/proxy/:instance", func(c echo.Context) error {
		callback := c.QueryParam("callback")
		if callback == "/rule/"+myRule.RuleID && c.Request().Method == http.MethodDelete {
			return c.String(http.StatusOK, "")
		}

		if callback == "/rule" {
			if c.Request().Method == http.MethodPost {
				return c.JSON(http.StatusOK, myRule)
			}

			return c.JSON(http.StatusOK, &acl.ServiceRuleData{
				ServiceInstance: types.ServiceInstance{
					BaseRules: []types.ServiceRule{
						myRule,
					},
				},
			})
		}
		t.Fatalf("method=%q, path=%q, callback=%q, err=\"Not found\"",
			c.Request().Method,
			c.Path(),
			callback,
		)
		return c.String(http.StatusNotFound, "")
	})

	fakeServer.HTTPErrorHandler = func(err error, c echo.Context) {
		t.Errorf("methods=%s, path=%s, err=%s", c.Request().Method, c.Path(), err.Error())
	}
	server := httptest.NewServer(fakeServer)
	os.Setenv("TSURU_TARGET", server.URL)

	resourceName := "acl_destination_rule.rule"
	config := `
	resource "acl_destination_rule" "rule" {
		instance =  "my-acl"
		app = "my-destination-app"
	}
					`

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:        config,
				ImportState:   true,
				ImportStateId: "acl-rule::acl::my-acl::app::my-destination-app",
				ResourceName:  "acl_destination_rule.rule",
			},
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app", "my-destination-app"),
					resource.TestCheckResourceAttr(resourceName, "service_name", "acl"),
					resource.TestCheckResourceAttr(resourceName, "instance", "my-acl"),
				),
			},
		},
	})
}

func TestAccResourceDestinationRuleDNS(t *testing.T) {
	fakeServer := echo.New()
	myRule := types.ServiceRule{
		Rule: types.Rule{
			RuleID: "my-rule",
			Destination: types.RuleType{
				ExternalDNS: &types.ExternalDNSRule{
					Name: "example.org",
					Ports: []types.ProtoPort{
						{
							Protocol: "TCP",
							Port:     80,
						},
						{
							Protocol: "TCP",
							Port:     443,
						},
					},
				},
			},
		},
	}

	fakeServer.Any("/services/acl/proxy/:instance", func(c echo.Context) error {
		callback := c.QueryParam("callback")
		if callback == "/rule/"+myRule.RuleID && c.Request().Method == http.MethodDelete {
			return c.String(http.StatusOK, "")
		}

		if callback == "/rule" {
			if c.Request().Method == http.MethodPost {
				return c.JSON(http.StatusOK, myRule)
			}

			return c.JSON(http.StatusOK, &acl.ServiceRuleData{
				ServiceInstance: types.ServiceInstance{
					BaseRules: []types.ServiceRule{
						myRule,
					},
				},
			})
		}
		t.Fatalf("method=%q, path=%q, callback=%q, err=\"Not found\"",
			c.Request().Method,
			c.Path(),
			callback,
		)
		return c.String(http.StatusNotFound, "")
	})

	fakeServer.HTTPErrorHandler = func(err error, c echo.Context) {
		t.Errorf("methods=%s, path=%s, err=%s", c.Request().Method, c.Path(), err.Error())
	}
	server := httptest.NewServer(fakeServer)
	os.Setenv("TSURU_TARGET", server.URL)

	resourceName := "acl_destination_rule.rule"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: `
resource "acl_destination_rule" "rule" {
	instance =  "my-acl"

	dns = "example.org"

	port {
		number   = 80
		protocol = "TCP"
	}

	port {
		number   = 443
		protocol = "TCP"
	}
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "dns", "example.org"),
					resource.TestCheckResourceAttr(resourceName, "port.0.number", "80"),
					resource.TestCheckResourceAttr(resourceName, "port.1.number", "443"),
				),
			},
		},
	})
}

func TestAccResourceDestinationRuleIP(t *testing.T) {
	fakeServer := echo.New()
	myRule := types.ServiceRule{
		Rule: types.Rule{
			RuleID: "my-rule",
			Destination: types.RuleType{
				ExternalIP: &types.ExternalIPRule{
					IP: "10.0.0.0/6",
					Ports: []types.ProtoPort{
						{
							Protocol: "TCP",
							Port:     80,
						},
						{
							Protocol: "TCP",
							Port:     443,
						},
					},
				},
			},
		},
	}

	fakeServer.Any("/services/acl/proxy/:instance", func(c echo.Context) error {
		callback := c.QueryParam("callback")
		if callback == "/rule/"+myRule.RuleID && c.Request().Method == http.MethodDelete {
			return c.String(http.StatusOK, "")
		}

		if callback == "/rule" {
			if c.Request().Method == http.MethodPost {
				return c.JSON(http.StatusOK, myRule)
			}

			return c.JSON(http.StatusOK, &acl.ServiceRuleData{
				ServiceInstance: types.ServiceInstance{
					BaseRules: []types.ServiceRule{
						myRule,
					},
				},
			})
		}
		t.Fatalf("method=%q, path=%q, callback=%q, err=\"Not found\"",
			c.Request().Method,
			c.Path(),
			callback,
		)
		return c.String(http.StatusNotFound, "")
	})

	fakeServer.HTTPErrorHandler = func(err error, c echo.Context) {
		t.Errorf("methods=%s, path=%s, err=%s", c.Request().Method, c.Path(), err.Error())
	}
	server := httptest.NewServer(fakeServer)
	os.Setenv("TSURU_TARGET", server.URL)

	resourceName := "acl_destination_rule.rule"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: `
resource "acl_destination_rule" "rule" {
	instance =  "my-acl"

	ip = "10.0.0.0/6"

	port {
		number   = 80
		protocol = "TCP"
	}

	port {
		number   = 443
		protocol = "TCP"
	}
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip", "10.0.0.0/6"),
					resource.TestCheckResourceAttr(resourceName, "port.0.number", "80"),
					resource.TestCheckResourceAttr(resourceName, "port.1.number", "443"),
				),
			},
		},
	})
}

func TestAccResourceDestinationRPaaS(t *testing.T) {
	fakeServer := echo.New()
	myRule := types.ServiceRule{
		Rule: types.Rule{
			RuleID: "my-rule",
			Destination: types.RuleType{
				RpaasInstance: &types.RpaasInstanceRule{
					ServiceName: "rpaasv2-be",
					Instance:    "my-rpaas",
				},
			},
		},
	}

	fakeServer.Any("/services/acl/proxy/:instance", func(c echo.Context) error {
		callback := c.QueryParam("callback")
		if callback == "/rule/"+myRule.RuleID && c.Request().Method == http.MethodDelete {
			return c.String(http.StatusOK, "")
		}

		if callback == "/rule" {
			if c.Request().Method == http.MethodPost {
				return c.JSON(http.StatusOK, myRule)
			}

			return c.JSON(http.StatusOK, &acl.ServiceRuleData{
				ServiceInstance: types.ServiceInstance{
					BaseRules: []types.ServiceRule{
						myRule,
					},
				},
			})
		}
		t.Fatalf("method=%q, path=%q, callback=%q, err=\"Not found\"",
			c.Request().Method,
			c.Path(),
			callback,
		)
		return c.String(http.StatusNotFound, "")
	})

	fakeServer.HTTPErrorHandler = func(err error, c echo.Context) {
		t.Errorf("methods=%s, path=%s, err=%s", c.Request().Method, c.Path(), err.Error())
	}
	server := httptest.NewServer(fakeServer)
	os.Setenv("TSURU_TARGET", server.URL)

	resourceName := "acl_destination_rule.rule"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: `
resource "acl_destination_rule" "rule" {
	instance =  "my-acl"

	rpaas {
		service_name = "rpaasv2-be"
		instance = "my-rpaas"
	} 
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rpaas.0.service_name", "rpaasv2-be"),
					resource.TestCheckResourceAttr(resourceName, "rpaas.0.instance", "my-rpaas"),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	tsuruTarget := os.Getenv("TSURU_TARGET")
	require.Contains(t, tsuruTarget, "http://127.0.0.1:")
}

func testAccResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		return nil
	}
}
