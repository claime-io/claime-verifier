package main

import (
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/contracts"
	"claime-verifier/lib/functions/lib/guild"
	guildrep "claime-verifier/lib/functions/lib/guild/persistence"
	"claime-verifier/lib/functions/lib/infrastructure/ethclient"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"claime-verifier/lib/functions/lib/transaction"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

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
	EOAInput struct {
		Signature string `json:"signature"`
		Message   string `json:"message"`
		RawTx     string `json:"rawTx"`
	}
	Input struct {
		Discord DiscordInput `json:"discord"`
		EOA     EOAInput     `json:"eoa"`
	}
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ssmClient := ssm.New(ctx)
	key, err := ssmClient.ClaimePublicKey()
	if err != nil {
		log.Error("get pubkey failed", err)
		return response(500), err
	}
	var in Input
	fmt.Println(request.Body)
	if err = json.Unmarshal([]byte(request.Body), &in); err != nil {
		log.Error("json unmarshal failed", err)
		return response(403), nil
	}
	if !verifyDiscordAppSignature(in.Discord, key) {
		log.Error("", errors.New("invalid signature"))
		return response(403), nil
	}
	if hasSignatureExpired(in.Discord) {
		log.Error("", errors.New("signature expired"))
		// TODO resend if expired
		return response(403), nil
	}
	rep := guildrep.New(ctx)
	guild, err := guild.New(ssmClient, rep)

	address, claim, err := recoverAddressAndClaim(in.EOA)
	if err != nil {
		log.Error("recover address failed", err)
		return response(400), nil
	}
	if claim.PropertyId != in.Discord.UserID {
		log.Error("", errors.New("invalid userID"))
		return response(403), nil
	}

	nfts, err := rep.ListContracts(in.Discord.GuildID)
	if err != nil {
		log.Error("", err)
		guild.ResendVerifyMessage(in.Discord.UserID, in.Discord.GuildID)
		return response(400), nil
	}
	granted := false

	if err != nil {
		log.Error("", err)
		return response(401), nil
	}

	for _, nft := range nfts {
		endpoint, err := ssmClient.EndpointByNetwork(nft.Network)
		if err != nil {
			log.Error("", err)
			continue
		}
		if isOwner(endpoint, common.HexToAddress(nft.ContractAddress), address) {
			if err = guild.GrantRole(in.Discord.UserID, nft); err != nil {
				log.Error("", err)
			}
			granted = true
		}
	}
	if granted {
		return response(200), nil
	}
	return response(401), nil
}

func isOwner(endpoint string, contractAddress common.Address, address common.Address) bool {
	cl, err := ethclient.NewERC721Client(endpoint)
	if err != nil {
		log.Error("", err)
		return false
	}
	caller, err := cl.Caller(contractAddress)
	if err != nil {
		log.Error("", err)
		return false
	}
	return caller.TokenOwner(address)
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

func recoverAddressAndClaim(in EOAInput) (common.Address, contracts.IClaimRegistryClaim, error) {
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

func verifyDiscordAppSignature(in DiscordInput, key ed25519.PublicKey) bool {
	vali, _ := strconv.ParseInt(in.Validity, 10, 64)
	return guild.Verify(guild.VerificationInput{
		SignatureInput: guild.SignatureInput{
			UserID:   in.UserID,
			GuildID:  in.GuildID,
			Validity: time.Unix(0, vali),
		},
		Sign: in.Signature,
	}, key)
}

func hasSignatureExpired(in DiscordInput) bool {
	vali, _ := strconv.ParseInt(in.Validity, 10, 64)
	return time.Now().After(time.Unix(0, vali))
}
