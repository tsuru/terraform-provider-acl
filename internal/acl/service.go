package acl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tsuru/acl-api/api/types"
)

// Create Rule
func (cli *clientImpl) DestinationRuleCreate(ctx context.Context, serviceName, instance string, rule *types.Rule) error {
	if len(serviceName) == 0 {
		return errors.New("Service Name not found")
	}
	if len(instance) == 0 {
		return errors.New("Service Instance not found")
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(rule)
	if err != nil {
		return err
	}

	rsp, err := doProxyRequest(ctx, http.MethodPost, serviceName, instance, "/rule", &buf, cli)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var savedRule types.Rule
	err = json.NewDecoder(rsp.Body).Decode(&savedRule)
	if err != nil {
		return err
	}

	if len(savedRule.RuleID) == 0 {
		return errors.New("Rule ID not found in response")
	}

	rule.RuleID = savedRule.RuleID
	return nil
}

// Get Rules
func (cli *clientImpl) DestinationRules(ctx context.Context, serviceName, instance string) (rules []types.Rule, err error) {
	if len(serviceName) == 0 {
		return nil, errors.New("Service Name not found")
	}
	if len(instance) == 0 {
		return nil, errors.New("Service Instance not found")
	}

	rsp, err := doProxyRequest(ctx, http.MethodGet, serviceName, instance, "/rule", nil, cli)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	ruleData := &ServiceRuleData{}
	err = json.NewDecoder(rsp.Body).Decode(ruleData)
	if err != nil {
		return
	}

	for _, rule := range ruleData.ServiceInstance.BaseRules {
		rules = append(rules, rule.Rule)
	}

	return
}

// Remove Rule
func (cli *clientImpl) DestinationRuleDelete(ctx context.Context, ruleID, serviceName, instance string) error {
	if len(serviceName) == 0 {
		serviceName = "acl"
	}
	if len(instance) == 0 {
		return errors.New("Service Instance not found")
	}

	rsp, err := doProxyRequest(ctx, http.MethodDelete, serviceName, instance, "/rule/"+ruleID, nil, cli)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	return nil
}

func ParseResourceID(id string) (*ParsedPrimaryID, error) {
	if len(id) == 0 {
		return nil, errors.New("ResourceId not found")
	}

	idParts := ParseIDParts(id)
	if len(idParts) < 3 {
		return nil, errors.New("Parse Resource ID invalid")
	}

	if len(idParts) == 3 {
		return &ParsedPrimaryID{
			Service:  idParts[0],
			Instance: idParts[1],
			RuleID:   idParts[2],
		}, nil
	}

	typeId := getID(0, idParts)
	if typeId != DestinationRulekey {
		return nil, errors.New("Destination Rule Key invalid")
	}

	parsedPrimaryID := &ParsedPrimaryID{}
	parsedPrimaryID.Service = getID(1, idParts)
	if len(parsedPrimaryID.Service) == 0 {
		return nil, errors.New("Service Name not found")
	}

	parsedPrimaryID.Instance = getID(2, idParts)
	if len(parsedPrimaryID.Instance) == 0 {
		return nil, errors.New("Service Instance not found")
	}

	parsedPrimaryID.Type = getID(3, idParts)
	if len(parsedPrimaryID.Type) == 0 {
		return nil, errors.New("Destination Type not found")
	}

	switch parsedPrimaryID.Type {
	case DestinationApp:
		parsedPrimaryID.AppName = getID(4, idParts)
	case DestinationPool:
		parsedPrimaryID.PoolName = getID(4, idParts)
	case DestinationCIDR:
		parsedPrimaryID.CIDR = getID(4, idParts)
	case DestinationDNS:
		parsedPrimaryID.DNS = getID(4, idParts)
	case DestinationRpaaS:
		parsedPrimaryID.RpaasService = getID(4, idParts)
		parsedPrimaryID.RpaasInstance = getID(5, idParts)
	default:
		return nil, errors.New("Parse Resource ID failed")
	}

	return parsedPrimaryID, nil
}
