package main

import (
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/contracts"
	"claime-verifier/lib/functions/lib/guild"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"claime-verifier/lib/functions/lib/transaction"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum/common"
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
		Message   string `json:"message"`
		RawTx     string `json:"rawTx"`
	}
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	k, err := ssm.New().ClaimePublicKey(ctx)
	if err != nil {
		log.Error("get pubkey failed", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
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
		// TODO resend if expired
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "{}",
			Headers:    lib.Headers(),
		}, errors.New("invalid signature")
	}

	address, claim, err := recoverAddressAndClaim(in)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "{}",
			Headers:    lib.Headers(),
		}, err
	}
	if claim.PropertyId != in.UserID {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "{}",
			Headers:    lib.Headers(),
		}, errors.New("invalid userID")
	}
	println(address)
	// TODO verify NFT ownership

	// TODO grant Role

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "{}",
		Headers:    lib.Headers(),
	}, nil
}

func main() {
	lambda.Start(handler)
}

func recoverAddressAndClaim(in Input) (common.Address, contracts.IClaimRegistryClaim, error) {
	if in.RawTx != "" {
		address, err := transaction.RecoverAddressFromTx(in.RawTx, in.Signature)
		if err != nil {
			return common.Address{}, contracts.IClaimRegistryClaim{}, err
		}
		claim, err := transaction.RecoverClaimFromTx(in.RawTx)
		if err != nil {
			return common.Address{}, contracts.IClaimRegistryClaim{}, err
		}
		return address, claim, nil
	}
	address, err := transaction.RecoverAddressFromMessage(in.Message, in.Signature)
	if err != nil {
		return common.Address{}, contracts.IClaimRegistryClaim{}, err
	}
	claim, err := transaction.RecoverClaimFromMessage(in.Message)
	if err != nil {
		return common.Address{}, contracts.IClaimRegistryClaim{}, err
	}
	return address, claim, nil
}
