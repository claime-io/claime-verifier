package main

import (
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/guild"
	guildrep "claime-verifier/lib/functions/lib/guild/persistence"
	"claime-verifier/lib/functions/lib/infrastructure/discord"
	slackclient "claime-verifier/lib/functions/lib/infrastructure/slack"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"errors"
	"fmt"

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
		Network         string `json:"network"`
	}
)

func handler(ctx context.Context, request map[string]interface{}) (interface{}, error) {
	keyresolver := ssm.New()
	fmt.Println(request)
	if !discord.VerifyInteractionRequest(ctx, request, keyresolver) {
		return unauthorized()
	}
	converter, err := discord.NewConverter(ctx, keyresolver)
	if err != nil {
		return unauthorized()
	}

	res, err := converter.HandleInteractionResponse(request)
	if err != nil {
		return unauthorized()
	}
	if !res.ShouldProcess() {
		return res, err
	}
	req, interaction, err := converter.ToRegisterContractInput(request)
	if err != nil {
		log.Error("", err)
		return unauthorized()
	}
	i, err := guild.New(ctx, keyresolver, guildrep.New())
	if err != nil {
		log.Error("", err)
		return unauthorized()
	}
	i.RegisterContract(ctx, interaction.ChannelID, req)
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
		Network:         d.Options[2].Value.(string),
	}
}
