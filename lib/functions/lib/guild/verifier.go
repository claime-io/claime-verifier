package guild

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"time"
)

type (
	SignatureInput struct {
		UserID, GuildID string
		Validity        time.Time
		Timestamp       time.Time
	}
	VerificationInput struct {
		SignatureInput
		Sign string
	}
)

func (in SignatureInput) String() string {
	return "timestamp=" + fmt.Sprint(in.Timestamp.UnixNano()) + "userId=" + in.UserID + "&guildId=" + in.GuildID + "&validity=" + fmt.Sprint(in.Validity.UnixNano())
}

func Sign(in SignatureInput, key ed25519.PrivateKey) string {
	return hex.EncodeToString(ed25519.Sign(key, []byte(in.String())))
}

func Verify(in VerificationInput, key ed25519.PublicKey) bool {
	return verify(in, key)
}

func verify(in VerificationInput, key ed25519.PublicKey) bool {
	sig, err := hex.DecodeString(in.Sign)
	if err != nil {
		return false
	}
	if len(sig) != ed25519.SignatureSize {
		return false
	}
	var msg bytes.Buffer
	msg.WriteString(in.SignatureInput.String())
	return ed25519.Verify(key, msg.Bytes(), sig)
}
