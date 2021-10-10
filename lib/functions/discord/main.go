package main

import (
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	guildrep "claime-verifier/lib/functions/lib/guild/persistence"
	"claime-verifier/lib/functions/lib/infrastructure/discord"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"errors"
	"fmt"

	"claime-verifier/lib/functions/lib/subscribe"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	SET_COMMAND_NAME    = "set"
	LIST_COMMAND_NAME   = "list"
	DELETE_COMMAND_NAME = "delete"
)

type (
	// RegisterContractInput register contract input
	RegisterContractInput struct {
		RoleID          string `json:"roleId"`
		ContractAddress string `json:"contract_address"`
		Network         string `json:"network"`
	}
)

func handler(ctx context.Context, request map[string]interface{}) (interface{}, error) {
	keyresolver := ssm.New(ctx)
	fmt.Println(request)
	if !discord.VerifyInteractionRequest(request, keyresolver) {
		return unauthorized()
	}
	res, err := discord.HandleInteractionResponse(request)
	if err != nil {
		log.Error("", err)
		return unauthorized()
	}
	if !res.ShouldProcess() {
		return res, err
	}
	interaction, err := discord.ToInteraction(request)
	if err != nil {
		log.Error("", err)
		return unauthorized()
	}
	interactor, err := guild.New(keyresolver, guildrep.New(ctx))
	if err != nil {
		log.Error("", err)
		return unauthorized()
	}
	if interaction.ApplicationCommandData().Name == SET_COMMAND_NAME {
		input := discord.ToRegisterContractInput(interaction.ApplicationCommandData(), interaction.GuildID)
		err = interactor.RegisterContract(ctx, interaction.ChannelID, interaction.GuildID, interaction.Member.Permissions, input)
		if err != nil {
			log.Error("RegisterContract", err)
		}
	}
	if interaction.ApplicationCommandData().Name == LIST_COMMAND_NAME {
		err = interactor.ListNFTs(ctx, interaction.ChannelID, interaction.GuildID)
		if err != nil {
			log.Error("List NFTs", err)
		}
	}
	if interaction.ApplicationCommandData().Name == DELETE_COMMAND_NAME {
		address := discord.ToDeleteContractAddressInput(interaction.ApplicationCommandData())
		err = interactor.DeleteNFT(ctx, interaction.ChannelID, interaction.GuildID, address)
		if err != nil {
			log.Error("Delete NFT", err)
		}
	}
	return res, err
}

func unauthorized() (interface{}, error) {
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
