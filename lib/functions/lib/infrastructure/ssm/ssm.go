package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/ethereum/go-ethereum/log"
)

type (
	// Client ssm client
	Client struct {
		svc ssmiface.SSMAPI
	}
	// Key key
	Key string
)

const (
	keyPrefix             = "claime-verifier-"
	infuraKeyPrefix       = keyPrefix + "infura-key-"
	slackTokenKey         = keyPrefix + "slack-token"
	slackSigningSecretKey = keyPrefix + "slack-signingsecret"
)

func keyOf(network string) string {
	return infuraKeyPrefix + network
}

// New New client
func New() Client {
	return Client{
		svc: ssm.New(session.New()),
	}
}

// WsEndpoint get ws endpoint
func (c Client) WsEndpoint(ctx context.Context, network string) (val string, err error) {
	return c.get(ctx, keyOf(network))
}

// SlackToken get slack token
func (c Client) SlackToken(ctx context.Context) (val string, err error) {
	return c.get(ctx, slackTokenKey)
}

// SlackSigningSecret get signing secret
func (c Client) SlackSigningSecret(ctx context.Context) (val string, err error) {
	return c.get(ctx, slackSigningSecretKey)
}

// Get get parameter
func (c Client) get(ctx context.Context, key string) (val string, err error) {
	out, err := c.svc.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Error("ssm get parameter failed", err)
		return
	}
	return *out.Parameter.Value, nil
}
