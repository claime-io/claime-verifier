package registry

import (
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
		if err != nil {
			t.Error(err)
		}
		if len(res) == 0 {
			t.Error("empty response")
		}
		claime := res[0]
		if claime.Evidence != "claime-dev.tk" {
			t.Error("got", claime.Evidence)
		}
		if claime.PropertyType != "Domain" {
			t.Error("got", claime.PropertyType)
		}
		if claime.PropertyID != "claime-dev.tk" {
			t.Error("got", claime.PropertyID)
		}
		if claime.Method != "TXT" {
			t.Error("got", claime.Method)
		}
	})
}
