package ssm

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"

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
	}
	// Key key
	Key string
)

const (
	keyPrefix             = "claime-verifier-"
	infuraKeyPrefix       = keyPrefix + "infura-key-"
	discordPublicKey      = keyPrefix + "discord-public-key"
	discordBotToken       = keyPrefix + "discord-bot-token"
	claimePublicKey       = keyPrefix + "public-key"
	claimePrivateKey      = keyPrefix + "private-key"
	endpointRinkeby       = keyPrefix + "endpoint-rinkeby"
	endpointMainnet       = keyPrefix + "endpoint-mainnet"
	endpointPolygon       = keyPrefix + "endpoint-polygon"
	twitterConsumerKey    = keyPrefix + "twitter-consumer-key"
	twitterConsumerSecret = keyPrefix + "twitter-consumer-secret"
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
	return c.get(ctx, withEnvSuffix(discordPublicKey))
}

func (c Client) EndpointByNetwork(ctx context.Context, network string) (val string, err error) {
	if network == "rinkeby" {
		return c.get(ctx, endpointRinkeby)
	}
	if network == "mainnet" {
		return c.get(ctx, endpointMainnet)
	}
	if network == "polygon" {
		return c.get(ctx, endpointPolygon)
	}
	return "", errors.New(fmt.Sprintf("Unsupported network : %s", network))
}

func (c Client) ClaimePublicKey(ctx context.Context) (val ed25519.PublicKey, err error) {
	return c.getKey(ctx, withEnvSuffix(claimePublicKey))
}

func (c Client) ClaimePrivateKey(ctx context.Context) (val ed25519.PrivateKey, err error) {
	return c.getKey(ctx, withEnvSuffix(claimePrivateKey))
}

func (c Client) TwitterConsumerKey(ctx context.Context) (val string, err error) {
	return c.get(ctx, withEnvSuffix(twitterConsumerKey))
}

func (c Client) TwitterConsumerSecret(ctx context.Context) (val string, err error) {
	return c.get(ctx, withEnvSuffix(twitterConsumerSecret))
}

func (c Client) getKey(ctx context.Context, key string) ([]byte, error) {
	k, err := c.get(ctx, key)
	if err != nil {
		return []byte{}, err
	}
	return hex.DecodeString(k)
}

func (c Client) DiscordBotToken(ctx context.Context) (val string, err error) {
	return c.get(ctx, withEnvSuffix(discordBotToken))
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

func withEnvSuffix(key string) string {
	env := os.Getenv("EnvironmentId")
	if env == "" {
		env = "dev"
	}
	if env == "prod" {
		return key
	}
	return key + "-" + env
}
