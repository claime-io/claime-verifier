package twitter

import (
	"claime-verifier/lib/functions/lib/common/log"
	"context"
	"strconv"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	evidencePrefix = "claime-ownership-claim="
)

type (
	// Client client
	Client struct {
		svc *twitter.Client
	}
	// Resolver resolver consumer key & secret
	Resolver interface {
		TwitterConsumerKey(ctx context.Context) (string, error)
		TwitterConsumerSecret(ctx context.Context) (string, error)
	}
)

// New new client
func New(ctx context.Context, r Resolver) (Client, error) {
	key, err := r.TwitterConsumerKey(ctx)
	if err != nil {
		return Client{}, err
	}
	sec, err := r.TwitterConsumerSecret(ctx)
	if err != nil {
		return Client{}, err
	}
	return new(key, sec), err
}

func new(cons, sec string) Client {
	config := &clientcredentials.Config{
		ClientID:     cons,
		ClientSecret: sec,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)
	return Client{
		svc: client,
	}
}

// EOA get eoa from twitter
func (c Client) EOA(ctx context.Context, id string) (common.Address, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error("id should be int64", err)
		return common.Address{}, err
	}
	ts, _, err := c.svc.Statuses.Lookup([]int64{i}, nil)
	if err != nil {
		log.Error("lookup tweet failed", err)
		return common.Address{}, err
	}
	if len(ts) == 0 {
		return common.Address{}, err
	}
	return eoa(ts[0].Text), nil
}

func eoa(rawMessage string) common.Address {
	return common.HexToAddress(eoaRaw(rawMessage))
}

func eoaRaw(raw string) string {
	exp := strings.TrimLeft(raw, evidencePrefix)
	return strings.TrimRight(strings.TrimLeft(exp, `"`), `"`)
}
