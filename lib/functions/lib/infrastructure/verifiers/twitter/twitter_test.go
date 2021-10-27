package twitter

import (
	"claime-verifier/lib/functions/lib/claim"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	validEvidence = "1448877989106651140"
	validClaim    = "claime-ownership-claim=0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"
)

var (
	verifyingEOA = common.HexToAddress("0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb")
)

type fakeLookUpper struct {
	twitterService
	FaceLookup func(ids []int64, params *twitter.StatusLookupParams) ([]twitter.Tweet, *http.Response, error)
}

func (mock fakeLookUpper) Lookup(ids []int64, params *twitter.StatusLookupParams) ([]twitter.Tweet, *http.Response, error) {
	return mock.FaceLookup(ids, params)
}

func newFakeLookUpper(tweet string, ID string, err error) fakeLookUpper {
	return fakeLookUpper{
		FaceLookup: func(ids []int64, params *twitter.StatusLookupParams) ([]twitter.Tweet, *http.Response, error) {
			return []twitter.Tweet{
				{
					Text: tweet,
					User: &twitter.User{
						IDStr: ID,
					},
				}}, nil, err
		},
	}
}

func TestEoaRaw(t *testing.T) {
	t.Run("get eoa success with "+validClaim, func(t *testing.T) {
		got := eoaRaw(validClaim)
		if common.HexToAddress(got).Hex() != verifyingEOA.Hex() {
			t.Error("got:", got)
		}
	})
}

func TestEoa(t *testing.T) {
	t.Run("get eoa success with "+validClaim, func(t *testing.T) {
		got := eoa(validClaim)
		assert.Equal(t, got, verifyingEOA)
	})
	invalidPrefixEvidence := strings.TrimPrefix(validClaim, "claime")
	t.Run("get eoa fail with invalid prefix"+invalidPrefixEvidence, func(t *testing.T) {
		got := eoa(invalidPrefixEvidence)
		assert.NotEqual(t, got, verifyingEOA)
	})
}

func TestEOA(t *testing.T) {
	propertyID := "claime"
	t.Run("get eoa", func(t *testing.T) {
		client := Client{
			lookupper: newFakeLookUpper(validClaim, propertyID, nil),
		}
		got, err := client.EOA(context.Background(), claim.Claim{Evidence: validEvidence})
		assert.Nil(t, err)
		assert.Equal(t, got.Actual, validClaim)
		assert.Equal(t, got.Got, verifyingEOA)
		assert.Equal(t, got.PropertyID, propertyID)
	})
	t.Run("error if evidence is not int64", func(t *testing.T) {
		client := Client{lookupper: newFakeLookUpper("", "", nil)}
		_, err := client.EOA(context.Background(), claim.Claim{Evidence: "string"})
		assert.Error(t, err)
	})
	t.Run("error if failed to lookup", func(t *testing.T) {
		client := Client{lookupper: newFakeLookUpper("", "", errors.Errorf(""))}
		_, err := client.EOA(context.Background(), claim.Claim{Evidence: validEvidence})
		assert.Error(t, err)
	})
	t.Run("error if tweet not found", func(t *testing.T) {
		client := Client{
			lookupper: fakeLookUpper{
				FaceLookup: func(ids []int64, params *twitter.StatusLookupParams) ([]twitter.Tweet, *http.Response, error) {
					return []twitter.Tweet{}, nil, nil
				},
			},
		}
		_, err := client.EOA(context.Background(), claim.Claim{Evidence: validEvidence})
		assert.Error(t, err)
	})
}
