package discord

import (
	"bytes"
	"claime-verifier/lib/functions/lib/common/log"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func VerifyRequest(request map[string]interface{}, publicKey string) bool {
	signature, ok := request["signature"].(string)
	if !ok {
		return false
	}
	timestamp, ok := request["timestamp"].(string)
	if !ok {
		return false
	}
	req, err := json.Marshal(request["jsonBody"])
	if err != nil {
		return false
	}

	httpreq, err := http.NewRequest("", "", bytes.NewReader(req))
	if err != nil {
		log.Error("", err)
		return false
	}
	key, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Error("", err)
		return false
	}
	httpreq.Header.Add("X-Signature-Ed25519", signature)
	httpreq.Header.Add("X-Signature-Timestamp", timestamp)
	return discordgo.VerifyInteraction(httpreq, key)
}
