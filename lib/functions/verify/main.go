package main

import (
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/claim"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/infrastructure/registry"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"claime-verifier/lib/functions/lib/transaction"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum/common"
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
	rep, err := registry.NewProvider(ctx, "rinkeby", ssmClient)
	if err != nil {
		log.Error("client initialize failed", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    lib.Headers(),
			Body:       "{}",
		}, nil
	}
	eoa := request.PathParameters["eoa"]
	address := common.HexToAddress(eoa)
	service := claim.NewService(rep)
	claims, err := service.Of(ctx, address)
	if err != nil {
		log.Error("get claim failed", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    lib.Headers(),
			Body:       "{}",
		}, nil
	}
	res, err := json.Marshal(&claims)
	if err != nil {
		log.Error("json marshal failed", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    lib.Headers(),
			Body:       "{}",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    lib.Headers(),
		Body:       string(res),
	}, err
}

func main() {
	lambda.Start(handler)
}
