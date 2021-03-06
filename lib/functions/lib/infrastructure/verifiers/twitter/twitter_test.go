package twitter

import (
	"claime-verifier/lib/functions/lib/claim"
	"context"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	validEvidence = "1448877989106651140"
	validClaim    = "claime-ownership-claim=0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb\nSee results in Claime:\nhttps://claime.io/0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"
)

var (
	verifyingEOA = common.HexToAddress("0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb")
)

type fakeLookUpper struct {
	twitterService
	fakeLookup func(id int64) (tweetEvidence, error)
}

func (mock fakeLookUpper) Lookup(id int64) (tweetEvidence, error) {
	return mock.fakeLookup(id)
}

func newFakeLookUpper(tweet string, userID string, err error) fakeLookUpper {
	return fakeLookUpper{
		fakeLookup: func(id int64) (tweetEvidence, error) {
			return tweetEvidence{
				text:   tweet,
				userID: userID,
			}, err
		},
	}
}

func TestEoaRaw(t *testing.T) {
	t.Run("get eoa success with "+validClaim, func(t *testing.T) {
		got := eoaRaw(validClaim)
		assert.Equal(t, verifyingEOA.Hex(), got)
	})
	t.Run("get empty with invalid claim", func(t *testing.T) {
		got := eoaRaw("xxxxx")
		assert.Equal(t, "", got)
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
		evidence, err := client.Find(context.Background(), claim.Claim{Evidence: validEvidence})
		assert.Nil(t, err)
		assert.Equal(t, evidence.PropertyID, propertyID)
		assert.Equal(t, verifyingEOA, evidence.EOAs[0])
		assert.Equal(t, validClaim, evidence.Evidences[0])
	})
	t.Run("error if evidence is not int64", func(t *testing.T) {
		client := Client{lookupper: newFakeLookUpper("", "", nil)}
		_, err := client.Find(context.Background(), claim.Claim{Evidence: "string"})
		assert.Error(t, err)
	})
	t.Run("error if failed to lookup", func(t *testing.T) {
		client := Client{lookupper: newFakeLookUpper("", "", errors.Errorf(""))}
		_, err := client.Find(context.Background(), claim.Claim{Evidence: validEvidence})
		assert.Error(t, err)
	})
}
