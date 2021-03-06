package discord

import (
	"bytes"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/ethereum/go-ethereum/common"
)

type (
	PubKeyResolver interface {
		DiscordPublicKey(ctx context.Context) (val string, err error)
	}
	KeyResolver interface {
		PubKeyResolver
		DiscordBotToken(ctx context.Context) (val string, err error)
		ClaimePrivateKey(ctx context.Context) (val ed25519.PrivateKey, err error)
	}
	InteractionConverter struct {
		pubkey    string
		bottoken  string
		claimeKey ed25519.PrivateKey
	}
	InteractionResponse struct {
		Type int `json:"type"`
	}
)

func (i InteractionResponse) ShouldProcess() bool {
	return i.Type == 5
}

func HandleInteractionResponse(request map[string]interface{}) (InteractionResponse, error) {
	req, err := json.Marshal(request["jsonBody"])
	if err != nil {
		return InteractionResponse{}, err
	}
	webhook := discordgo.Webhook{}
	err = json.Unmarshal(req, &webhook)
	if err != nil {
		return InteractionResponse{}, err
	}
	if webhook.Type == discordgo.WebhookTypeChannelFollower {
		return InteractionResponse{Type: 5}, nil
	}
	if webhook.Type == discordgo.WebhookTypeIncoming {
		return InteractionResponse{Type: 1}, nil
	}
	return InteractionResponse{}, errors.New("unknown type:" + strconv.Itoa(int(webhook.Type)))
}

func ToInteraction(request map[string]interface{}) (discordgo.Interaction, error) {
	req, err := json.Marshal(request["jsonBody"])
	if err != nil {
		log.Error("", err)
		return discordgo.Interaction{}, err
	}
	interaction := discordgo.Interaction{}
	if err = interaction.UnmarshalJSON(req); err != nil {
		log.Error("", err)
		return discordgo.Interaction{}, err
	}
	return interaction, nil
}

// VerifyInteractionRequest verify signature of interaction request
func VerifyInteractionRequest(ctx context.Context, request map[string]interface{}, r PubKeyResolver) bool {
	pubkey, err := r.DiscordPublicKey(ctx)
	if err != nil {
		log.Error("", err)
		return false
	}
	return verify(request, pubkey)
}

func verify(request map[string]interface{}, publicKey string) bool {
	signature, ok := request["signature"].(string)
	if !ok {
		log.Info("sig not found")
		return false
	}
	timestamp, ok := request["timestamp"].(string)
	if !ok {
		log.Info("timestamp not found")
		return false
	}
	req, err := json.Marshal(request["jsonBody"])
	if err != nil {
		log.Error("", err)
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

func ToRegisterContractInput(d discordgo.ApplicationCommandInteractionData, guildID string) guild.NFTInfo {
	if len(d.Options) < 3 {
		return guild.NFTInfo{}
	}
	return guild.NFTInfo{
		RoleID:          d.Options[0].Value.(string),
		ContractAddress: d.Options[1].Value.(string),
		Network:         d.Options[2].Value.(string),
		GuildID:         guildID,
	}
}

func ToDeleteContractAddressInput(d discordgo.ApplicationCommandInteractionData) common.Address {
	return common.HexToAddress(d.Options[0].Value.(string))
}
