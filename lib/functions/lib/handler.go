package lib

import (
	"claime-verifier/lib/functions/lib/claim"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"claime-verifier/lib/functions/lib/infrastructure/verifiers/twitter"
	"claime-verifier/lib/functions/lib/infrastructure/verifiers/txt"
	"claime-verifier/lib/functions/lib/infrastructure/verifiers/website"
	"context"
)

// Headers with headers
func Headers() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Methods":     "GET",
		"Access-Control-Allow-Credentials": "false",
		"Access-Control-Allow-Origin":      "*",
	}
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
