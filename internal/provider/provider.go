package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tsuru/terraform-provider-acl/internal/acl"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Description: "Target to tsuru API",
				Optional:    true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "Token to authenticate on tsuru API (optional)",
				Optional:    true,
			},
			"skip_cert_verification": {
				Type:        schema.TypeBool,
				Description: "Disable certificate verification",
				Default:     false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TSURU_SKIP_CERT_VERIFICATION", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"acl_destination_rule": resourceACLDestinationRule(),
		},
	}
	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, p.TerraformVersion)
	}
	return p
}

type aclProvider struct {
	client           acl.Client
	terraformVersion string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	p := &aclProvider{}

	host := d.Get("host").(string)
	token := d.Get("token").(string)

	cli, err := acl.NewClient(ctx, host, token)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	p.client = cli
	return p, diags
}
