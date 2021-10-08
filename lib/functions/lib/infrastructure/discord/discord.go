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
	Gateway             struct{}
	InteractionResponse struct {
		Type int `json:"type"`
	}
)

func (i InteractionResponse) ShouldProcess() bool {
	return i.Type == 2
}

func NewConverter(ctx context.Context, r KeyResolver) (InteractionConverter, error) {
	pub, err := r.DiscordPublicKey(ctx)
	if err != nil {
		return InteractionConverter{}, err
	}
	t, err := r.DiscordBotToken(ctx)
	if err != nil {
		return InteractionConverter{}, err
	}
	pri, err := r.ClaimePrivateKey(ctx)
	if err != nil {
		return InteractionConverter{}, err
	}
	return InteractionConverter{
		pubkey:    pub,
		bottoken:  t,
		claimeKey: pri,
	}, nil
}

func (s InteractionConverter) HandleInteractionResponse(request map[string]interface{}) (InteractionResponse, error) {
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
		return InteractionResponse{Type: 1}, nil
	}
	if webhook.Type == discordgo.WebhookTypeIncoming {
		return InteractionResponse{Type: 4}, nil
	}
	return InteractionResponse{}, errors.New("unknown type:" + strconv.Itoa(int(webhook.Type)))
}

func toInteraction(request map[string]interface{}) (discordgo.Interaction, error) {
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

func (s InteractionConverter) ToRegisterContractInput(request map[string]interface{}) (guild.RegisterContractInput, error) {
	interaction, err := toInteraction(request)
	if err != nil {
		return guild.RegisterContractInput{}, err
	}
	return toInput(interaction.ApplicationCommandData()), err
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

func toInput(d discordgo.ApplicationCommandInteractionData) guild.RegisterContractInput {
	if len(d.Options) < 3 {
		return guild.RegisterContractInput{}
	}
	return guild.RegisterContractInput{
		RoleID:          d.Options[0].Value.(string),
		ContractAddress: d.Options[1].Value.(string),
		Network:         d.Options[2].Value.(string),
	}
}
