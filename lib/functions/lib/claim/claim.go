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
		Claim
		Verified bool      `json:"verified"`
		Actual   string    `json:"actual"`
		At       time.Time `json:"at"`
	}

	// Verifier verifier
	Verifier struct {
		PropertyType, Method string
		Default              bool
	}

	// EOAOutput eoa
	EOAOutput struct {
		Actual     string
		Got        common.Address
		PropertyID string
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
)

// NewService new service
func NewService(rep Repository, supported map[Verifier]EvidenceRepository) Service {
	return Service{
		rep:       rep,
		verifiers: supported,
	}
}

// VerifiedClaims list verified claims of eoa.
func (s Service) VerifiedClaims(ctx context.Context, eoa common.Address) ([]VerifiedOutput, error) {
	claims, err := s.claimsOf(ctx, eoa)
	if err != nil {
		return []VerifiedOutput{}, err
	}
	res := []VerifiedOutput{}
	for _, cl := range claims {
		verifier, ok := supportedVerifier(cl, s.verifiers)
		if !ok {
			continue
		}
		got, err := s.verifiers[verifier].EOA(ctx, cl)
		if err != nil {
			continue
		}
		res = append(res, VerifiedOutput{
			Claim:    cl,
			Actual:   got.Actual,
			At:       time.Now(),
			Verified: verified(cl, eoa, got),
		})
	}
	return res, nil
}

func supportedVerifier(c Claim, verifiers map[Verifier]EvidenceRepository) (Verifier, bool) {
	for k := range verifiers {
		if k.PropertyType != c.PropertyType {
			continue
		}
		if c.Method == "" && k.Default {
			return k, true
		}
		if c.Method == k.Method {
			return k, true
		}
	}
	return Verifier{}, false
}

func verified(cl Claim, eoa common.Address, got EOAOutput) bool {
	return (cl.PropertyID == got.PropertyID) && (eoa == got.Got)
}

func (s Service) claimsOf(ctx context.Context, eoa common.Address) ([]Claim, error) {
	return s.rep.ClaimsOf(ctx, eoa)
}
