package acl

import (
	"strings"

	"github.com/tsuru/acl-api/api/types"
)

func GenerateID(values []string) string {
	var newValues []string
	for _, v := range values {
		if len(v) > 0 {
			newValues = append(newValues, v)
		}
	}

	newId := strings.Join(newValues, "::")
	return strings.TrimSpace(newId)
}

func ParseIDParts(id string) []string {
	if len(id) == 0 {
		return []string{}
	}

	return strings.Split(strings.TrimSpace(id), "::")
}

func getID(key int, ids []string) string {
	if len(ids) <= key {
		return ""
	}

	return ids[key]
}

func FindRuleByParsedPrimaryID(rules []types.Rule, id *ParsedPrimaryID) (rule *types.Rule) {
	for _, v := range rules {
		// Rule ID
		if len(id.RuleID) > 0 && v.RuleID == id.RuleID {
			return &v
		}

		// App Name
		if v.Destination.TsuruApp != nil && len(id.AppName) > 0 {
			if v.Destination.TsuruApp.AppName == id.AppName {
				return &v
			}
		}

		// Pool Name
		if v.Destination.TsuruApp != nil && len(id.PoolName) > 0 {
			if v.Destination.TsuruApp.PoolName == id.PoolName {
				return &v
			}
		}

		// Rpaas Instance
		if v.Destination.RpaasInstance != nil && len(id.RpaasService) > 0 && len(id.RpaasInstance) > 0 {
			if v.Destination.RpaasInstance.ServiceName == id.RpaasService {
				if v.Destination.RpaasInstance.Instance == id.RpaasInstance {
					return &v
				}
			}
		}

		// CIDR / IP
		if v.Destination.ExternalIP != nil && len(id.CIDR) > 0 {
			if v.Destination.ExternalIP.IP == id.CIDR {
				return &v
			}
		}

		// DNS
		if v.Destination.ExternalDNS != nil && len(id.DNS) > 0 {
			if v.Destination.ExternalDNS.Name == id.DNS {
				return &v
			}
		}
	}

	return nil
}

func FindRuleBySingleID(rules []types.Rule, id string) (rule *types.Rule) {
	for _, v := range rules {
		if v.RuleID == id {
			return &v
		}
	}

	return nil
}
