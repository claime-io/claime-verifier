package ssm

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"

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
		ctx context.Context
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
func New(ctx context.Context) Client {
	return Client{
		svc: ssm.New(session.New(&aws.Config{Region: aws.String("us-east-1")})),
		ctx: ctx,
	}
}

// WsEndpoint get ws endpoint
func (c Client) WsEndpoint(network string) (val string, err error) {
	return c.get(keyOf(network))
}

// DiscordPublicKey get Discord public key
func (c Client) DiscordPublicKey() (val string, err error) {
	return c.get(discordPublicKey)
}

func (c Client) EndpointByNetwork(network string) (val string, err error) {
	if network == "rinkeby" {
		return c.get(endpointRinkeby)
	}
	if network == "mainnet" {
		return c.get(endpointMainnet)
	}
	return "", errors.New(fmt.Sprintf("Unsupported network : %s", network))
}

func (c Client) ClaimePublicKey() (val ed25519.PublicKey, err error) {
	return c.getKey(claimePublicKey)
}

func (c Client) ClaimePrivateKey() (val ed25519.PrivateKey, err error) {
	return c.getKey(claimePrivateKey)
}

func (c Client) getKey(key string) ([]byte, error) {
	k, err := c.get(key)
	if err != nil {
		return []byte{}, err
	}
	return hex.DecodeString(k)
}

func (c Client) DiscordBotToken() (val string, err error) {
	return c.get(discordBotToken)
}

// SlackToken get slack token
func (c Client) SlackToken() (val string, err error) {
	return c.get(slackTokenKey)
}

// SlackSigningSecret get signing secret
func (c Client) SlackSigningSecret() (val string, err error) {
	return c.get(slackSigningSecretKey)
}

// Get get parameter
func (c Client) get(key string) (val string, err error) {
	out, err := c.svc.GetParameterWithContext(c.ctx, &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Error("ssm get parameter failed", err)
		return
	}

	return *out.Parameter.Value, nil
}
