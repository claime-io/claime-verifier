package lib

import (
	"claime-verifier/lib/functions/lib/common/log"
	slackclient "claime-verifier/lib/functions/lib/infrastructure/slack"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// Headers with headers
func Headers(request events.APIGatewayProxyRequest) map[string]string {
	return map[string]string{}
}

// NewParameterStore new parameter store
func NewParameterStore(ctx context.Context) (ssm.Client, error) {
	return ssm.New(), nil
}

// Verify verify slack message
func Verify(request events.APIGatewayProxyRequest, cli slackclient.Client) error {
	return cli.Verify(request)
}

// Slackcli new slack client
func Slackcli(ctx context.Context, parameterstore ssm.Client) (slackclient.Client, error) {
	token, err := parameterstore.SlackToken(ctx)
	if err != nil {
		log.Error("retrive token failed", err)
	}
	signingsecret, err := parameterstore.SlackSigningSecret(ctx)
	if err != nil {
		log.Error("retrive sign secret failed", err)
	}
	return slackclient.New(token, signingsecret), nil
}
