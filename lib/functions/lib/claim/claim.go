package claim

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type (
	// Claim claim
	Claim struct {
		PropertyType string `json:"propertyType"`
		PropertyID   string `json:"propertyId"`
		Evidence     string `json:"evidence"`
		Method       string `json:"method"`
	}

	// VerifiedOutput output
	VerifiedOutput struct {
		Claim  Claim              `json:"claim"`
		Result VerificationResult `json:"result"`
		Actual Actual             `json:"actual"`
		At     time.Time          `json:"at"`
		Error  string             `json:"error"`
	}

	// Actual actual
	Actual struct {
		PropertyID string `json:"propertyId"`
		Evidence   string `json:"evidence"`
	}

	// Verifier verifier
	Verifier struct {
		PropertyType, Method string
	}

	// EOAOutput eoa
	EOAOutput struct {
		Actual Actual
		EOA    common.Address
	}

	// Repository repository
	Repository interface {
		ClaimsOf(ctx context.Context, eoa common.Address) ([]Claim, error)
	}

	// Service service
	Service struct {
		rep       Repository
		verifiers map[Verifier]EvidenceRepository
	}
	// EvidenceRepository evidence repository
	EvidenceRepository interface {
		EOA(ctx context.Context, cl Claim) (EOAOutput, error)
	}
	VerificationResult string
)

const (
	verified    = VerificationResult("Verified")
	failed      = VerificationResult("Failed")
	unsupported = VerificationResult("Unsupported")
)

// NewService new service
func NewService(rep Repository, supported map[Verifier]EvidenceRepository) Service {
	return Service{
		rep:       rep,
		verifiers: supported,
	}
}

// VerifyClaims lists the verification results for claims of eoa.
func (s Service) VerifyClaims(ctx context.Context, eoa common.Address) ([]VerifiedOutput, error) {
	claims, err := s.claimsOf(ctx, eoa)
	if err != nil {
		return []VerifiedOutput{}, err
	}
	res := []VerifiedOutput{}
	for _, cl := range claims {
		verifier := s.verifiers[Verifier{PropertyType: cl.PropertyType, Method: cl.Method}]
		if verifier == nil {
			res = append(res, VerifiedOutput{
				Claim:  cl,
				At:     time.Now(),
				Result: unsupported,
			})
			continue
		}
		got, err := verifier.EOA(ctx, cl)
		if err != nil {
			res = append(res, VerifiedOutput{
				Claim:  cl,
				At:     time.Now(),
				Result: failed,
				Error:  err.Error(),
			})
			continue
		}
		res = append(res, VerifiedOutput{
			Claim:  cl,
			Actual: got.Actual,
			At:     time.Now(),
			Result: verify(cl, eoa, got),
		})
	}
	return res, nil
}

func verify(cl Claim, eoa common.Address, got EOAOutput) VerificationResult {
	if (cl.PropertyID == got.Actual.PropertyID) && (eoa == got.EOA) {
		return verified
	}
	return failed
}

func (s Service) claimsOf(ctx context.Context, eoa common.Address) ([]Claim, error) {
	return s.rep.ClaimsOf(ctx, eoa)
}
