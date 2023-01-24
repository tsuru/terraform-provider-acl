package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tsuru/acl-api/api/types"
)

func isRetryableError(err error) bool {
	return strings.Contains(err.Error(), "event locked")
}

func ruleFromResource(d *schema.ResourceData) *types.Rule {
	rule := &types.Rule{}

	ports := d.Get("port").([]interface{})
	var protoPorts []types.ProtoPort
	for _, port := range ports {
		portMap := port.(map[string]interface{})
		proto := portMap["protocol"].(string)
		port := portMap["number"].(int)
		protoPorts = append(protoPorts, types.ProtoPort{
			Protocol: proto,
			Port:     uint16(port),
		})
	}

	rule.Destination.TsuruApp = parseTsuruApp(d)
	rule.Destination.RpaasInstance = parseRpaas(d)

	dns := d.Get("dns").(string)
	if dns != "" {
		rule.Destination.ExternalDNS = &types.ExternalDNSRule{
			Name:  dns,
			Ports: protoPorts,
		}
	}
	dstIP := d.Get("ip").(string)
	if dstIP != "" {
		rule.Destination.ExternalIP = &types.ExternalIPRule{
			IP:    dstIP,
			Ports: protoPorts,
		}
	}

	return rule
}

func parseRpaas(d *schema.ResourceData) *types.RpaasInstanceRule {
	list := d.Get("rpaas").([]interface{})
	if len(list) == 0 {
		return nil
	}
	sourceRpaas := list[0].(map[string]interface{})
	return &types.RpaasInstanceRule{
		ServiceName: sourceRpaas["service_name"].(string),
		Instance:    sourceRpaas["instance"].(string),
	}
}

func flattenRpaas(rpaas *types.RpaasInstanceRule) []interface{} {
	if rpaas == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"service_name": rpaas.ServiceName,
			"instance":     rpaas.Instance,
		},
	}
}

func flattenProtoPorts(ports []types.ProtoPort) []interface{} {
	if ports == nil {
		return nil
	}

	var portList []interface{}
	for _, port := range ports {
		portList = append(portList, map[string]interface{}{
			"protocol": port.Protocol,
			"number":   int(port.Port),
		})
	}

	return portList
}

func parseTsuruApp(d *schema.ResourceData) *types.TsuruAppRule {
	app, _ := d.Get("app").(string)
	pool, _ := d.Get("pool").(string)
	if app == "" && pool == "" {
		return nil
	}

	if app != "" {
		return &types.TsuruAppRule{
			AppName: app,
		}
	}

	return &types.TsuruAppRule{
		PoolName: pool,
	}
}

func tsuruSchema(baseName string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{baseName + ".0.pool"},
			},
			"pool": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{baseName + ".0.name"},
			},
		},
	}
}

func rpaasSchema(baseName string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{baseName + ".0.service_name"},
			},
			"instance": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{baseName + ".0.instance"},
			},
		},
	}
}
