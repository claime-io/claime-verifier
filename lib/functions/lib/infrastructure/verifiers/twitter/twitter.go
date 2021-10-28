package twitter

import (
	"claime-verifier/lib/functions/lib/claim"
	"claime-verifier/lib/functions/lib/common/log"
	"context"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const (
	evidencePrefix = "claime-ownership-claim="
)

type (
	// Client client
	Client struct {
		lookupper tweetLookUpper
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
	return Client{
		lookupper: newTwitterService(key, sec),
	}, err
}

// EOA get eoa from twitter
func (c Client) EOA(ctx context.Context, cl claim.Claim) (claim.EOAOutput, error) {
	i, err := strconv.ParseInt(cl.Evidence, 10, 64)
	if err != nil {
		log.Error("id should be int64", err)
		return claim.EOAOutput{}, err
	}
	tweet, err := c.lookupper.Lookup(i)
	if err != nil {
		return claim.EOAOutput{}, err
	}
	return claim.EOAOutput{
		Actual: claim.Actual{
			Evidence:   tweet.text,
			PropertyID: tweet.userID,
		},
		EOA: eoa(tweet.text),
	}, nil
}

func eoa(rawMessage string) common.Address {
	return common.HexToAddress(eoaRaw(rawMessage))
}

func eoaRaw(raw string) string {
	return strings.TrimPrefix(raw, evidencePrefix)
}
