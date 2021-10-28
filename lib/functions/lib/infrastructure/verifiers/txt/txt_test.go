package txt

import (
	"claime-verifier/lib/functions/lib/claim"
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEOA(t *testing.T) {
	t.Run("enable to get eoa", func(t *testing.T) {
		want := common.HexToAddress("0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb")
		got, err := Client{}.EOA(context.Background(), claim.Claim{PropertyID: "claime-dev.tk"})
		if err != nil {
			t.Error(err)
		}
		if want != got.EOA {
			t.Error("got: ", got.EOA.String())
		}
	})
}
