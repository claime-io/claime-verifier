package registry

import (
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestClaimsOf(t *testing.T) {
	ctx := context.Background()
	resolver := ssm.New()
	t.Run("enable to list claims", func(t *testing.T) {
		p, err := NewProvider(ctx, "rinkeby", resolver)
		if err != nil {
			t.Error(err)
		}
		res, err := p.ClaimsOf(ctx, common.HexToAddress("0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"))
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, len(res), 1)
		assert.Equal(t, res[0].PropertyType, "Domain")
		assert.Equal(t, res[0].PropertyID, "claime-dev.tk")
		assert.Equal(t, res[0].Method, "TXT")
		assert.Equal(t, res[0].Evidence, "")
		assert.Equal(t, res[0].Network, "rinkeby")
	})
}
