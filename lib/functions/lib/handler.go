package lib

import (
	"claime-verifier/lib/functions/lib/claim"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"claime-verifier/lib/functions/lib/infrastructure/verifiers/twitter"
	"claime-verifier/lib/functions/lib/infrastructure/verifiers/txt"
	"context"
	"os"
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

// SupportedVerifiers list supported verifiers by Verify Methods
func SupportedVerifiers(ctx context.Context, ssm ssm.Client) (map[claim.Verifier]claim.EvidenceRepository, error) {
	tw, err := twitter.New(ctx, ssm)
	return map[claim.Verifier]claim.EvidenceRepository{{PropertyType: "Domain", Method: "TXT", Default: true}: txt.New(), {PropertyType: "Twitter Account", Method: "Tweet", Default: true}: tw}, err
}
