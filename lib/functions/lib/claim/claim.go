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

	// VerificationResponse output
	VerificationResponse struct {
		Claim    Claim              `json:"claim"`
		Result   VerificationResult `json:"result"`
		Evidence Evidence           `json:"actual"`
		At       time.Time          `json:"at"`
		Error    string             `json:"error"`
	}

	// Evidence actual
	Evidence struct {
		PropertyID string           `json:"propertyId"`
		EOAs       []common.Address `json:"eoas"`
		Evidences  []string         `json:"evidences"`
	}

	// VerificationKey verificationKey
	PropertyKey struct {
		PropertyType, Method string
	}

	// EvidenceRepository evidenceRepository
	EvidenceRepository interface {
		Find(ctx context.Context, cl Claim) (Evidence, error)
	}

	// Repository repository
	Repository interface {
		ClaimsOf(ctx context.Context, eoa common.Address) ([]Claim, error)
	}

	// Service service
	Service struct {
		rep          Repository
		repositories map[PropertyKey]EvidenceRepository
	}
	VerificationResult string
)

const (
	verified    = VerificationResult("Verified")
	failed      = VerificationResult("Failed")
	unsupported = VerificationResult("Unsupported")
)

// NewService new service
func NewService(rep Repository, supported map[PropertyKey]EvidenceRepository) Service {
	return Service{
		rep:          rep,
		repositories: supported,
	}
}

// VerifyClaims lists the verification results for claims of eoa.
func (s Service) VerifyClaims(ctx context.Context, eoa common.Address) ([]VerificationResponse, error) {
	claims, err := s.claimsOf(ctx, eoa)
	if err != nil {
		return []VerificationResponse{}, err
	}
	res := []VerificationResponse{}
	for _, cl := range claims {
		repository := s.repositories[PropertyKey{PropertyType: cl.PropertyType, Method: cl.Method}]
		if repository == nil {
			res = append(res, VerificationResponse{
				Claim:  cl,
				At:     time.Now(),
				Result: unsupported,
			})
			continue
		}
		evidence, err := repository.Find(ctx, cl)
		if err != nil {
			res = append(res, VerificationResponse{
				Claim:  cl,
				At:     time.Now(),
				Result: failed,
				Error:  err.Error(),
			})
			continue
		}
		res = append(res, VerificationResponse{
			Claim:    cl,
			Evidence: evidence,
			Result:   verify(cl, eoa, evidence),
			At:       time.Now(),
		})
	}
	return res, nil
}

func (s Service) claimsOf(ctx context.Context, eoa common.Address) ([]Claim, error) {
	return s.rep.ClaimsOf(ctx, eoa)
}

func verify(cl Claim, eoa common.Address, evidence Evidence) VerificationResult {
	if cl.PropertyID != evidence.PropertyID {
		return failed
	}
	for _, actualEOA := range evidence.EOAs {
		if eoa == actualEOA {
			return verified
		}
	}
	return failed
}
