package twitter

import (
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
)

type (
	twitterService struct {
		svc *twitter.Client
	}
	tweetLookUpper interface {
		Lookup(ids []int64, params *twitter.StatusLookupParams) ([]twitter.Tweet, *http.Response, error)
	}
)

func (t twitterService) Lookup(ids []int64, params *twitter.StatusLookupParams) ([]twitter.Tweet, *http.Response, error) {
	return t.svc.Statuses.Lookup(ids, params)
}
