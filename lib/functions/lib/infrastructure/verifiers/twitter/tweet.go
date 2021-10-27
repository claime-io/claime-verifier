package twitter

import (
	"claime-verifier/lib/functions/lib/common/log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	tokenURL = "https://api.twitter.com/oauth2/token"
)

type (
	twitterService struct {
		svc *twitter.Client
	}
	tweetLookUpper interface {
		Lookup(id int64) (tweetEvidence, error)
	}
	tweetEvidence struct {
		text   string
		userID string
	}
)

func newTwitterService(cons, sec string) twitterService {
	config := &clientcredentials.Config{
		ClientID:     cons,
		ClientSecret: sec,
		TokenURL:     tokenURL,
	}
	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)
	return twitterService{
		svc: client,
	}
}

func (t twitterService) Lookup(id int64) (tweetEvidence, error) {
	tweets, _, err := t.svc.Statuses.Lookup([]int64{id}, nil)
	if err != nil {
		log.Error("failed to lookup tweet.", err)
		return tweetEvidence{}, err
	}
	return toEvidence(tweets, id)
}

func toEvidence(tweets []twitter.Tweet, id int64) (tweetEvidence, error) {
	if len(tweets) == 0 {
		return tweetEvidence{}, errors.Errorf("Tweet not found: %d", id)
	}
	return tweetEvidence{
		text:   tweets[0].Text,
		userID: tweets[0].User.ScreenName,
	}, nil
}
