package acl

import (
	"github.com/tsuru/acl-api/api/types"
)

const (
	DestinationApp   = "app"
	DestinationPool  = "pool"
	DestinationRpaaS = "rpaas"
	DestinationCIDR  = "ip"
	DestinationDNS   = "dns"

	DestinationRulekey = "acl-rule"
)

var Destinations = []string{
	DestinationApp,
	DestinationPool,
	DestinationRpaaS,
	DestinationCIDR,
	DestinationDNS,
}

type ServiceRuleData struct {
	ServiceInstance types.ServiceInstance
}

type ParsedPrimaryID struct {
	Service  string
	Instance string

	RuleID string

	Type string

	AppName  string
	PoolName string
	CIDR     string
	DNS      string

	RpaasService  string
	RpaasInstance string
}
