package main

import (
	"bytes"
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	slackclient "claime-verifier/lib/functions/lib/infrastructure/slack"
	"context"
	"encoding/hex"
	"encoding/json"
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
)

func handler(ctx context.Context, request map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	token := request["signature"].(string)
	timestamp := request["timestamp"].(string)

	req, _ := json.Marshal(request["jsonBody"])
	webhook := discordgo.Webhook{}
	json.Unmarshal(req, &webhook)
	httpreq, err := http.NewRequest("", "", bytes.NewReader(req))
	if err != nil {
		log.Error("", err)
	}
	key, _ := hex.DecodeString("API_KEY")
	if err != nil {
		log.Error("", err)
	}
	httpreq.Header.Add("X-Signature-Ed25519", token)
	httpreq.Header.Add("X-Signature-Timestamp", timestamp)
	result := discordgo.VerifyInteraction(httpreq, key)
	fmt.Println(result)
	if !result {
		return events.APIGatewayProxyResponse{
			StatusCode:      401,
			Body:            `[UNAUTHORIZED] invalid request signature`,
			Headers:         map[string]string{},
			IsBase64Encoded: false,
		}, nil
	}
	fmt.Println("type")
	fmt.Println(webhook.Type)

	if webhook.Type == discordgo.WebhookTypeChannelFollower {
		resp := struct {
			Type int `json:"type"`
		}{
			Type: 1,
		}
		rp, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{
			StatusCode:      200,
			Body:            string(rp),
			Headers:         map[string]string{},
			IsBase64Encoded: false,
		}, nil
	}
	if webhook.Type == discordgo.WebhookTypeIncoming {
		resp := struct {
			Type int `json:"type"`
		}{
			Type: 5,
		}
		rp, _ := json.Marshal(resp)
		return events.APIGatewayProxyResponse{
			StatusCode:      200,
			Body:            string(rp),
			Headers:         map[string]string{},
			IsBase64Encoded: false,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      401,
		Body:            `[UNAUTHORIZED] invalid request signature`,
		Headers:         map[string]string{},
		IsBase64Encoded: false,
	}, nil
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
		Headers:    lib.Headers(request),
		Body:       string(out),
	}, nil
}

func main() {
	lambda.Start(handler)
}

func newApp(ctx context.Context, slcli slackclient.Client) subscribe.Registrar {
	return subscribe.NewRegistrar(slcli, repository.New())
}
