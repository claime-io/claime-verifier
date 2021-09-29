package main

import (
	"claime-verifier/lib/functions/lib"
	slackclient "claime-verifier/lib/functions/lib/infrastructure/slack"
	"context"
	"encoding/json"

	"claime-verifier/lib/functions/lib/subscribe"
	repository "claime-verifier/lib/functions/lib/subscribe/persistence"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	requiredArgs = 3
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ps, err := lib.NewParameterStore(ctx)
	if err != nil {
		return unexpectedError(request, err)
	}
	s, err := lib.Slackcli(ctx, ps)
	if err != nil {
		return unexpectedError(request, err)
	}
	if err = s.Verify(request); err != nil {
		return unexpectedError(request, err)
	}
	app := newApp(ctx, s)
	res, err := app.Register(ctx, request.Body)
	if err != nil {
		return errorResponse(request, res)
	}
	resJSON, _ := json.Marshal(&res)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    lib.Headers(request),
		Body:       string(resJSON),
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
