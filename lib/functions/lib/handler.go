package lib

import (
	"claime-verifier/lib/functions/lib/claim"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"claime-verifier/lib/functions/lib/infrastructure/verifiers/twitter"
	"claime-verifier/lib/functions/lib/infrastructure/verifiers/txt"
	"claime-verifier/lib/functions/lib/infrastructure/verifiers/website"
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// Headers with headers
func Headers(origin string) map[string]string {
	return map[string]string{
		"Access-Control-Allow-Headers":     "*",
		"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Origin":      allowedOrigin(origin),
	}
}

func allowedOrigin(origin string) string {
	env := os.Getenv("EnvironmentId")
	if env == "prod" {
		return os.Getenv("AllowedOrigin")
	}
	return origin
}

// Origin retrive origin value from request headers
func Origin(request events.APIGatewayProxyRequest) string {
	return request.Headers["origin"]
}

// SupportedVerifications list supported verifications by Property and verification methods
func SupportedVerifications(ctx context.Context, ssm ssm.Client) (map[claim.PropertyKey]claim.EvidenceRepository, error) {
	tw, err := twitter.New(ctx, ssm)
	return map[claim.PropertyKey]claim.EvidenceRepository{
		{PropertyType: "Domain", Method: "TXT"}:            txt.New(),
		{PropertyType: "Twitter Account", Method: "Tweet"}: tw,
		{PropertyType: "Website", Method: "Meta Tag"}:      website.New(),
	}, err
}
