package acl

import (
	"context"

	"github.com/tsuru/acl-api/api/types"
	tsuruCmd "github.com/tsuru/tsuru/cmd"
)

type Client interface {
	DestinationRuleCreate(ctx context.Context, serviceName, instance string, rule *types.Rule) error
	DestinationRules(ctx context.Context, serviceName, instance string) (rules []types.Rule, err error)
	DestinationRuleDelete(ctx context.Context, ruleID, serviceName, instance string) error
}

type clientImpl struct {
	Host  string
	token string
}

func NewClient(ctx context.Context, host, token string) (Client, error) {
	if len(host) == 0 {
		target, err := tsuruCmd.GetTarget()
		if err != nil {
			return nil, err
		}
		host = target
	}

	if len(token) == 0 {
		tsuruToken, err := tsuruCmd.ReadToken()
		if err != nil {
			return nil, err
		}
		token = tsuruToken
	}

	return &clientImpl{
		Host:  host,
		token: token,
	}, nil
}
