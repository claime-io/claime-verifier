package subgraph

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type fakeClient struct {
	subgraphClient
	fakeClaims func(ctx context.Context, eoa common.Address) (listClaims, error)
}

func (mock fakeClient) claims(ctx context.Context, eoa common.Address) (listClaims, error) {
	return mock.fakeClaims(ctx, eoa)
}

func newFakeClient(claims listClaims, err error) fakeClient {
	return fakeClient{
		fakeClaims: func(ctx context.Context, eoa common.Address) (listClaims, error) {
			return claims, err
		},
	}
}

const (
	url          = "https://api.studio.thegraph.com/query/8417/claime-rinkeby/v0.0.1"
	addressStr   = "0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"
	propertyType = "Twitter"
	propertyID   = "example"
	method       = "Tweet"
	evidence     = "1234567890"
	network      = "mainnet"
)

var (
	address    = common.HexToAddress(addressStr)
	mockClaims = listClaims{
		Claims: []claimData{
			{
				PropertyType: propertyType,
				PropertyID:   propertyID,
				Method:       method,
				Evidence:     evidence,
				Network:      network,
			},
		},
	}
)

func TestClaimsOf(t *testing.T) {
	t.Run("can mapping claims", func(t *testing.T) {
		client := newFromClient(newFakeClient(mockClaims, nil))
		claims, err := client.ClaimsOf(context.Background(), address)
		assert.Nil(t, err)
		assert.Len(t, claims, 1)
		assert.Equal(t, claims[0].PropertyType, propertyType)
		assert.Equal(t, claims[0].PropertyID, propertyID)
		assert.Equal(t, claims[0].Method, method)
		assert.Equal(t, claims[0].Evidence, evidence)
		assert.Equal(t, claims[0].Network, network)
	})
	t.Run("can return err", func(t *testing.T) {
		errMsg := "error"
		client := newFromClient(newFakeClient(listClaims{}, errors.New(errMsg)))
		_, err := client.ClaimsOf(context.Background(), address)
		assert.Equal(t, err.Error(), errMsg)

	})
	t.Run("can call subgraph", func(t *testing.T) {
		t.Skip()
		client := New(url)
		claims, err := client.ClaimsOf(context.Background(), address)
		assert.Nil(t, err)
		assert.Len(t, claims, 1)
		assert.Equal(t, claims[0].PropertyType, "Domain")
		assert.Equal(t, claims[0].PropertyID, "claime-dev.tk")
		assert.Equal(t, claims[0].Method, "TXT")
		assert.Equal(t, claims[0].Evidence, "")
		assert.Equal(t, claims[0].Network, "rinkeby")
	})
}
