package twitter

import (
	"claime-verifier/lib/functions/lib/claim"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"testing"
)

func TestEoaRaw(t *testing.T) {
	// TODO: pass tests
	texts := []string{`claime-ownership-claim="0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"`, `this is test claime-ownership-claim="0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb" test`, `claime-ownership-claim="0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"\n`}
	for _, text := range texts {
		t.Run("get eoa success with "+text, func(t *testing.T) {
			got := eoaRaw(text)
			if got != "0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb" {
				t.Error("got:", got)
			}
		})
	}
}

func TestEoa(t *testing.T) {
	// TODO: pass tests

	texts := []string{`claime-ownership-claim=\"0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb\"`, `this is test claime-ownership-claim="0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb" test`, `claime-ownership-claim="0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"\n`}
	for _, text := range texts {
		t.Run("get eoa success with "+text, func(t *testing.T) {
			got := eoa(text)
			if got.String() != "0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb" {
				t.Error("got:", got.String())
			}
		})
	}
}

func TestEOA(t *testing.T) {
	t.Run("get eoa", func(t *testing.T) {
		resolver := ssm.New()
		target, err := New(context.Background(), resolver)
		if err != nil {
			t.Error(err)
		}
		got, err := target.EOA(context.Background(), claim.Claim{Evidence: "1448877989106651140"})
		if err != nil {
			t.Error(err)
		}
		if got.Got.String() != "0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb" {
			t.Error("got:", got.Got.String())
		}
	})
}
