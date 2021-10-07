package main

import (
	"bytes"
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	slackclient "claime-verifier/lib/functions/lib/infrastructure/slack"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"claime-verifier/lib/functions/lib/subscribe"
	repository "claime-verifier/lib/functions/lib/subscribe/persistence"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

const (
	requiredArgs = 3
	mockGuildID  = "892441777808765049"
)

var (
	token string
)

type (
	// RegisterContractInput register contract input
	RegisterContractInput struct {
		RoleID          string `json:"roleId"`
		ContractAddress string `json:"contract_address"`
		ChainID         int    `json:"chain_id"`
	}
)

func handler(ctx context.Context, request map[string]interface{}) (interface{}, error) {
	token := request["signature"].(string)
	timestamp := request["timestamp"].(string)

	token, err := ssm.New().DiscordBotToken(context.Background())
	if err != nil {
		log.Error("error get bot token", err)
		return "error get bot token", nil
	}

	req, _ := json.Marshal(request["jsonBody"])
	webhook := discordgo.Webhook{}
	interaction := discordgo.Interaction{}
	interaction.UnmarshalJSON(req)
	fmt.Printf("%+v\n", request["jsonBody"])
	fmt.Printf("%+v\n", interaction)
	json.Unmarshal(req, &webhook)
	fmt.Printf("%+v\n", webhook)
	if guild.HasPermissionAdministrator(interaction.Member.Permissions) {
		fmt.Printf("called by admin")
	}

	dg, err := discordgo.New("Bot " + token)
	guildID := mockGuildID
	member, err := dg.GuildMember(guildID, interaction.User.ID)
	fmt.Printf("%+v\n", member)

	httpreq, err := http.NewRequest("", "", bytes.NewReader(req))
	if err != nil {
		log.Error("", err)
	}
	k, err := ssm.New().DiscordPublicKey(ctx)
	if err != nil {
		log.Error("", err)
	}
	key, _ := hex.DecodeString(k)
	if err != nil {
		log.Error("", err)
	}
	fmt.Println("key")
	fmt.Println(string(key))
	httpreq.Header.Add("X-Signature-Ed25519", token)
	httpreq.Header.Add("X-Signature-Timestamp", timestamp)
	result := discordgo.VerifyInteraction(httpreq, key)
	fmt.Println(result)
	if !result {
		return `[UNAUTHORIZED] invalid request signature`, errors.New("[UNAUTHORIZED] invalid request signature")
	}
	fmt.Println("type")
	fmt.Println(webhook.Type)

	if webhook.Type == discordgo.WebhookTypeChannelFollower {
		resp := struct {
			Type int `json:"type"`
		}{
			Type: 1,
		}
		return resp, nil
	}
	if webhook.Type == discordgo.WebhookTypeIncoming {

		resp := struct {
			Type int `json:"type"`
		}{
			Type: 4,
		}
		return resp, nil
	}
	return `[UNAUTHORIZED] invalid request signature`, errors.New("[UNAUTHORIZED] invalid request signature")

}

func unexpectedError(request events.APIGatewayProxyRequest, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}, err
}

func errorResponse(request events.APIGatewayProxyRequest, output subscribe.Output) (events.APIGatewayProxyResponse, error) {
	out, err := output.Parse()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    lib.Headers(),
		Body:       string(out),
	}, nil
}

func main() {
	lambda.Start(handler)
}

func newApp(ctx context.Context, slcli slackclient.Client) subscribe.Registrar {
	return subscribe.NewRegistrar(slcli, repository.New())
}

func toInput(d discordgo.ApplicationCommandInteractionData) RegisterContractInput {
	if len(d.Options) < 3 {
		return RegisterContractInput{}
	}
	return RegisterContractInput{
		RoleID:          d.Options[0].Value.(string),
		ContractAddress: d.Options[1].Value.(string),
		ChainID:         d.Options[2].Value.(int),
	}
}
