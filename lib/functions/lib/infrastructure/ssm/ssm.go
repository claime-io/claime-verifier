package ssm

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"

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
	discordPublicKey      = keyPrefix + "discord-public-key"
	discordBotToken       = keyPrefix + "discord-bot-token"
	claimePublicKey       = keyPrefix + "public-key"
	claimePrivateKey      = keyPrefix + "private-key"
	endpointRinkeby       = keyPrefix + "endpoint-rinkeby"
	endpointMainnet       = keyPrefix + "endpoint-mainnet"
)

func keyOf(network string) string {
	return infuraKeyPrefix + network
}

// New New client
func New() Client {
	return Client{
		svc: ssm.New(session.New(&aws.Config{Region: aws.String("us-east-1")})),
	}
}

// WsEndpoint get ws endpoint
func (c Client) WsEndpoint(ctx context.Context, network string) (val string, err error) {
	return c.get(ctx, keyOf(network))
}

// DiscordPublicKey get Discord public key
func (c Client) DiscordPublicKey(ctx context.Context) (val string, err error) {
	return c.get(ctx, discordPublicKey)
}

func (c Client) EndpointRinkeby(ctx context.Context) (val string, err error) {
	return c.get(ctx, endpointRinkeby)
}

func (c Client) EndpointMainnet(ctx context.Context) (val string, err error) {
	return c.get(ctx, endpointMainnet)
}

func (c Client) ClaimePublicKey(ctx context.Context) (val ed25519.PublicKey, err error) {
	return c.getKey(ctx, claimePublicKey)
}

func (c Client) ClaimePrivateKey(ctx context.Context) (val ed25519.PrivateKey, err error) {
	return c.getKey(ctx, claimePrivateKey)
}

func (c Client) getKey(ctx context.Context, key string) ([]byte, error) {
	k, err := c.get(ctx, key)
	if err != nil {
		return []byte{}, err
	}
	return hex.DecodeString(k)
}

func (c Client) DiscordBotToken(ctx context.Context) (val string, err error) {
	return c.get(ctx, discordBotToken)
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
