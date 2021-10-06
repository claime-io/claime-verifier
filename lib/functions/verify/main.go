package main

import (
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	requiredArgs = 3
)

type (
	Input struct {
		UserID    string `json:"userId"`
		GuildID   string `json:"guildId"`
		Validity  int64  `json:"validity"`
		Timestamp int64  `json:"timestamp"`
		Signature string `json:"signature"`
	}
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	k, err := ssm.New().ClaimePublicKey(ctx)
	if err != nil {
		log.Error("get pubkey failed", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "{}",
			Headers:    lib.Headers(),
		}, err
	}
	var in Input
	if err = json.Unmarshal([]byte(request.Body), &in); err != nil {
		log.Error("json unmarshal failed", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "{}",
			Headers:    lib.Headers(),
		}, err
	}
	if !guild.Verify(guild.VerificationInput{
		SignatureInput: guild.SignatureInput{
			UserID:    in.UserID,
			GuildID:   in.GuildID,
			Validity:  time.Unix(0, in.Validity),
			Timestamp: time.Unix(0, in.Timestamp),
		},
	}, k) {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "{}",
			Headers:    lib.Headers(),
		}, errors.New("invalid signature")
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "{}",
		Headers:    lib.Headers(),
	}, nil
}

func main() {
	lambda.Start(handler)
}
