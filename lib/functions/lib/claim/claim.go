package claim

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type (
	// Claim claim
	Claim struct {
		PropertyType string `json:"propertyType"`
		PropertyId   string `json:"propertyId"`
		Evidence     string `json:"evidence"`
		Method       string `json:"method"`
	}

	PropertyType string
	Method       string

	Repository interface {
		ClaimsOf(ctx context.Context, eoa common.Address) ([]Claim, error)
	}
	Service struct {
		rep       Repository
		verifiers map[string]EvidenceRepository
	}
	EvidenceRepository interface {
		EOA(ctx context.Context, propertyID string) (common.Address, error)
	}
)

var (
	SupportedPropertyTypes = []PropertyType{"Domain"}
	SupportedMethods       = []Method{"TXT", "Tweet"}
)

func NewService(rep Repository, supported map[string]EvidenceRepository) Service {
	return Service{
		rep:       rep,
		verifiers: supported,
	}
}

func (s Service) Of(ctx context.Context, eoa common.Address) ([]Claim, error) {
	return s.rep.ClaimsOf(ctx, eoa)
}

// VerifiedClaims list verified claims of eoa.
func (s Service) VerifiedClaims(ctx context.Context, eoa common.Address) ([]Claim, error) {
	claims, err := s.claimsOf(ctx, eoa)
	if err != nil {
		return []Claim{}, err
	}
	res := []Claim{}
	for _, cl := range claims {
		m := cl.Method
		verifier, ok := s.verifiers[m]
		if !ok {
			continue
		}
		got, err := verifier.EOA(ctx, cl.PropertyId)
		if err != nil {
			continue
		}
		if verified(eoa, got) {
			res = append(res, cl)
		}
	}
	return res, nil
}

func verified(want, got common.Address) bool {
	return want.String() == got.String()
}

func (s Service) claimsOf(ctx context.Context, eoa common.Address) ([]Claim, error) {
	return s.rep.ClaimsOf(ctx, eoa)
}
