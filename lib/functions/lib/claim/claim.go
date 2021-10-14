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
		rep Repository
	}
)

var (
	SupportedPropertyTypes = []PropertyType{"Domain"}
	SupportedMethods       = []Method{"TXT"}
)

func NewService(rep Repository) Service {
	return Service{
		rep: rep,
	}
}

func (s Service) Of(ctx context.Context, eoa common.Address) ([]Claim, error) {
	return s.rep.ClaimsOf(ctx, eoa)
}
