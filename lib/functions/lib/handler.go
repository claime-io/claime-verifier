package lib

import (
	"claime-verifier/lib/functions/lib/common/log"
	slackclient "claime-verifier/lib/functions/lib/infrastructure/slack"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// Headers with headers
func Headers() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Headers":     "*",
		"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Origin":      os.Getenv("AllowedOrigin"),
	}
}

// NewParameterStore new parameter store
func NewParameterStore(ctx context.Context) (ssm.Client, error) {
	return ssm.New(ctx), nil
}

// Verify verify slack message
func Verify(request events.APIGatewayProxyRequest, cli slackclient.Client) error {
	return cli.Verify(request)
}

// Slackcli new slack client
func Slackcli(parameterstore ssm.Client) (slackclient.Client, error) {
	token, err := parameterstore.SlackToken()
	if err != nil {
		log.Error("retrive token failed", err)
	}
	signingsecret, err := parameterstore.SlackSigningSecret()
	if err != nil {
		log.Error("retrive sign secret failed", err)
	}
	return slackclient.New(token, signingsecret), nil
}
