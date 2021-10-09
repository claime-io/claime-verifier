package guild

import (
	"crypto/ed25519"
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	pub, pri, err := ed25519.GenerateKey(nil)
	input := SignatureInput{
		UserID:   "test",
		GuildID:  "guild",
		Validity: time.Now(),
	}
	signature := Sign(input, pri)
	t.Run("signature verified", func(t *testing.T) {
		if err != nil {
			t.Error("")
		}

		if !verify(VerificationInput{
			SignatureInput: input,
			Sign:           signature,
		}, pub) {
			t.Error("signature doesnt match")
		}
	})
	t.Run("signature not verified with fake userid", func(t *testing.T) {
		if verify(VerificationInput{
			SignatureInput: SignatureInput{
				UserID:   input.UserID + "a",
				GuildID:  input.GuildID,
				Validity: input.Validity,
			},
			Sign: signature,
		}, pub) {
			t.Error("signature matched")
		}
	})
	t.Run("signature not verified with fake guild id", func(t *testing.T) {
		if verify(VerificationInput{
			SignatureInput: SignatureInput{
				UserID:   input.UserID,
				GuildID:  input.GuildID + "a",
				Validity: input.Validity,
			},
			Sign: signature,
		}, pub) {
			t.Error("signature matched")
		}
	})
	t.Run("signature not verified with fake validity", func(t *testing.T) {
		if verify(VerificationInput{
			SignatureInput: SignatureInput{
				UserID:   input.UserID,
				GuildID:  input.GuildID,
				Validity: input.Validity.Add(1),
			},
			Sign: signature,
		}, pub) {
			t.Error("signature matched")
		}
	})
}
