package twitter

import (
	"claime-verifier/lib/functions/lib/common/log"
	"context"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type (
	Client struct {
		svc *twitter.Client
	}
	Resolver interface {
		TwitterConsumerKey(ctx context.Context) (string, error)
		TwitterComsumerSecret(ctx context.Context) (string, error)
	}
)

func New(ctx context.Context, r Resolver) (Client, error) {
	key, err := r.TwitterConsumerKey(ctx)
	if err != nil {
		return Client{}, err
	}
	sec, err := r.TwitterComsumerSecret(ctx)
	if err != nil {
		return Client{}, err
	}
	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: sec,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)
	return Client{
		svc: client,
	}, nil
}

func (c Client) Get(ctx context.Context, id int64) (twitter.Tweet, error) {
	ts, _, err := c.svc.Statuses.Lookup([]int64{id}, nil)
	if err != nil {
		log.Error("lookup tweet failed", err)
		return twitter.Tweet{}, err
	}
	if len(ts) == 0 {
		return twitter.Tweet{}, err
	}
	return ts[0], nil
}
