package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tsuru/acl-api/api/types"
	"github.com/tsuru/terraform-provider-acl/internal/acl"
)

func resourceACLDestinationRule() *schema.Resource {
	oneDestination := acl.Destinations

	return &schema.Resource{
		CreateContext: resourceACLDestinationRuleCreate,
		ReadContext:   resourceACLDestinationRuleRead,
		DeleteContext: resourceACLDestinationRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceACLDestinationRuleImport,
		},

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL Instance Name",
			},
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "acl",
				Description: "ACL Service Name",
			},

			"ip": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ExactlyOneOf: oneDestination,
				ValidateFunc: validation.IsCIDR,
				Description:  "Destination IP address",
			},

			"dns": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ExactlyOneOf: oneDestination,
				Description:  "Destination fully qualified domain name (FQDN)",
			},

			"app": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ExactlyOneOf: oneDestination,
				Description:  "Destination tsuru app name",
			},

			"pool": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ExactlyOneOf: oneDestination,
				Description:  "Tsuru Pool name",
			},

			"rpaas": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeList,
				MaxItems:     1,
				MinItems:     1,
				ExactlyOneOf: oneDestination,
				Elem:         rpaasSchema("rpaas"),
				Description:  "Destination tsuru rpaas name",
			},

			"port": {
				Optional:      true,
				ForceNew:      true,
				Type:          schema.TypeList,
				Description:   "Destination port and protocol list",
				ConflictsWith: []string{"app", "rpaas"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Description:  "Procotol name (ex: TCP, UDP, tcp, udp...)",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "tcp", "udp"}, false),
						},
						"number": {
							Type:         schema.TypeInt,
							Description:  "Port number",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsPortNumber,
						},
					},
				},
			},
		},
	}
}

func resourceACLDestinationRuleImport(ctx context.Context, rd *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	primaryID, err := acl.ParseResourceID(rd.Id())

	if err != nil {
		return nil, err
	}

	cli := m.(*aclProvider).client

	rules, err := cli.DestinationRules(ctx, primaryID.Service, primaryID.Instance)

	if err != nil {
		return nil, err
	}

	rule := acl.FindRuleByParsedPrimaryID(rules, primaryID)
	if rule == nil {
		return nil, errors.New("rule not found")
	}

	rd.Set("service_name", primaryID.Service)
	rd.Set("instance", primaryID.Instance)
	rd.SetId(rule.RuleID)

	return []*schema.ResourceData{rd}, nil
}

func resourceACLDestinationRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	cli := m.(*aclProvider).client

	serviceName := d.Get("service_name").(string)
	instance := d.Get("instance").(string)
	rule := ruleFromResource(d)

	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err := cli.DestinationRuleCreate(ctx, serviceName, instance, rule)
		if err != nil {
			if isRetryableError(err) {
				return resource.RetryableError(err)
			}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "DestinationRuleCreate",
				Detail:   err.Error(),
			})

			return resource.NonRetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		d.SetId(rule.RuleID)
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceACLDestinationRuleRead(ctx, d, m)
}

func resourceACLDestinationRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cli := m.(*aclProvider).client

	rule, err := readRuleFromResourceData(ctx, cli, d)
	if err != nil {
		return diag.FromErr(err)
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	// Destination TsuruApp
	if rule.Destination.TsuruApp == nil {
		if err := d.Set("app", nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("pool", nil); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("app", rule.Destination.TsuruApp.AppName); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("pool", rule.Destination.TsuruApp.PoolName); err != nil {
			return diag.FromErr(err)
		}
	}

	// Destination IP
	if rule.Destination.ExternalIP == nil {
		if err := d.Set("ip", nil); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("ip", rule.Destination.ExternalIP.IP); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("port", flattenProtoPorts(rule.Destination.ExternalIP.Ports)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Destination DNS
	if rule.Destination.ExternalDNS == nil {
		if err := d.Set("dns", nil); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("dns", rule.Destination.ExternalDNS.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("port", flattenProtoPorts(rule.Destination.ExternalDNS.Ports)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Destination Rpaas Instance
	if rule.Destination.RpaasInstance == nil {
		if err := d.Set("rpaas", nil); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("rpaas", flattenRpaas(rule.Destination.RpaasInstance)); err != nil {
			return diag.FromErr(err)
		}
	}

	if rule.Destination.ExternalIP == nil && rule.Destination.ExternalDNS == nil {
		if err := d.Set("port", nil); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceACLDestinationRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	cli := m.(*aclProvider).client

	rule, err := readRuleFromResourceData(ctx, cli, d)
	if err != nil {
		return diag.FromErr(err)
	}

	if rule == nil {
		return nil
	}

	serviceName := d.Get("service_name").(string)
	instance := d.Get("instance").(string)

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err := cli.DestinationRuleDelete(ctx, rule.RuleID, serviceName, instance)
		if err != nil {
			if isRetryableError(err) {
				return resource.RetryableError(err)
			}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "DestinationRuleDelete",
				Detail:   err.Error(),
			})

			return resource.NonRetryableError(err)
		}

		d.SetId("")
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readRuleFromResourceData(ctx context.Context, cli acl.Client, d *schema.ResourceData) (rule *types.Rule, err error) {
	serviceName := d.Get("service_name").(string)
	instance := d.Get("instance").(string)

	parts := acl.ParseIDParts(d.Id())
	fullID := d.Id()
	id := fullID

	if len(parts) != 1 && len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID: %s", fullID)
	}

	if len(parts) == 3 {
		serviceName = parts[0]
		instance = parts[1]
		id = parts[2]
	}

	rules, err := cli.DestinationRules(ctx, serviceName, instance)
	if err != nil {
		return nil, err
	}

	rule = acl.FindRuleBySingleID(rules, id)
	return rule, nil
}
