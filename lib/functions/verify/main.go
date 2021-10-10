package main

import (
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	guildrep "claime-verifier/lib/functions/lib/guild/persistence"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"claime-verifier/lib/functions/lib/transaction"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type (
	DiscordInput struct {
		UserID    string `json:"userId"`
		GuildID   string `json:"guildId"`
		Validity  string `json:"validity"`
		Signature string `json:"signature"`
	}

	Input struct {
		Discord DiscordInput         `json:"discord"`
		EOA     transaction.EOAInput `json:"eoa"`
	}
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ssmClient := ssm.New()
	var in guild.GrantRoleInput
	if err := json.Unmarshal([]byte(request.Body), &in); err != nil {
		log.Error("json unmarshal failed", err)
		return response(403), nil
	}
	rep := guildrep.New()
	guild, err := guild.New(ctx, ssmClient, rep)
	if err != nil {
		log.Error("", err)
		return response(500), nil
	}
	out, err := guild.Grant(ctx, in)
	return handleResponse(out, err)
}

func handleResponse(out guild.GrantRoleOutput, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return response(400), nil
	}
	if !out.ValidSig {
		return response(403), nil
	}
	if out.Expired {
		return response(403), nil
	}
	if out.Granted {
		return response(200), nil
	}
	return response(401), nil
}

func main() {
	lambda.Start(handler)
}

func response(statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       "{}",
		Headers:    lib.Headers(),
	}
}
