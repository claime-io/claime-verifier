package claim

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
)

type fakeClaimRepository struct {
	Repository
	fakeClaimsOf func(ctx context.Context, eoa common.Address) ([]Claim, error)
}

func (mock fakeClaimRepository) ClaimsOf(ctx context.Context, eoa common.Address) ([]Claim, error) {
	return mock.fakeClaimsOf(ctx, eoa)
}

func newFakeClaimRepository(claims []Claim, err error) fakeClaimRepository {
	return fakeClaimRepository{
		fakeClaimsOf: func(ctx context.Context, eoa common.Address) ([]Claim, error) {
			return claims, err
		},
	}
}

type fakeEvidenceRepository struct {
	EvidenceRepository
	fakeOutput func(ctx context.Context, cl Claim) (Evidence, error)
}

func (fake fakeEvidenceRepository) Find(ctx context.Context, cl Claim) (Evidence, error) {
	return fake.fakeOutput(ctx, cl)
}

func newFakeEvidenceRepository(out Evidence, err error) EvidenceRepository {
	return fakeEvidenceRepository{
		fakeOutput: func(ctx context.Context, cl Claim) (Evidence, error) {
			return out, err
		},
	}
}

func newFakeEvidenceRepositories(key PropertyKey, out Evidence, err error) map[PropertyKey]EvidenceRepository {
	return map[PropertyKey]EvidenceRepository{
		key: newFakeEvidenceRepository(out, err),
	}
}

var (
	verifyingEOA = common.HexToAddress("0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb")
	mockClaim    = Claim{
		PropertyType: "TestPropertyType",
		PropertyID:   "TestPropertyID",
		Evidence:     "",
		Method:       "TestMethod",
		Network:      "TestNetwork",
	}
	mockPropertyKey = PropertyKey{
		PropertyType: mockClaim.PropertyType,
		Method:       mockClaim.Method,
	}
	mockEvidence = Evidence{
		PropertyID: mockClaim.PropertyID,
		Evidences:  []string{"TestEvidence"},
		EOAs:       []common.Address{verifyingEOA},
	}
)

func TestVerifyClaim(t *testing.T) {
	t.Run("return verified if matched actual EOA", func(t *testing.T) {
		assert.Equal(t, verify(mockClaim, verifyingEOA, mockEvidence), verified)
	})
	t.Run("return failed if not matched got EOA", func(t *testing.T) {
		assert.Equal(t, verify(mockClaim, common.HexToAddress("anotherEOA"), mockEvidence), failed)
	})
	t.Run("return failed if not matched got propertyID", func(t *testing.T) {
		assert.Equal(t, verify(mockClaim, common.HexToAddress("anotherEOA"), Evidence{
			PropertyID: "another Property",
			EOAs:       mockEvidence.EOAs,
			Evidences:  mockEvidence.Evidences,
		}), failed)
	})
}
func TestVerifyClaims(t *testing.T) {
	t.Run("return verified if verified", func(t *testing.T) {
		service := NewService(
			newFakeClaimRepository([]Claim{mockClaim}, nil),
			newFakeEvidenceRepositories(mockPropertyKey, mockEvidence, nil),
		)
		outputs, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Nil(t, err)
		assert.NotEmpty(t, outputs)
		assert.Equal(t, outputs[0].Result, verified)
		assert.Equal(t, outputs[0].Evidence, mockEvidence)
		assert.Equal(t, outputs[0].Claim, mockClaim)
		assert.Equal(t, "", outputs[0].Error)
	})
	t.Run("return failed and message if error", func(t *testing.T) {
		expectedMessage := "Error Message"
		service := NewService(
			newFakeClaimRepository([]Claim{mockClaim}, nil),
			newFakeEvidenceRepositories(mockPropertyKey, Evidence{}, errors.Errorf(expectedMessage)),
		)
		outputs, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Nil(t, err)
		assert.NotEmpty(t, outputs)
		assert.Equal(t, outputs[0].Result, failed)
		assert.Equal(t, outputs[0].Error, expectedMessage)
	})
	t.Run("return unsuported if verifier not found", func(t *testing.T) {
		service := NewService(
			newFakeClaimRepository([]Claim{mockClaim}, nil),
			newFakeEvidenceRepositories(PropertyKey{}, Evidence{}, nil),
		)
		outputs, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Nil(t, err)
		assert.NotEmpty(t, outputs)
		assert.Equal(t, outputs[0].Result, unsupported)
	})
	t.Run("return error if failed to get claims", func(t *testing.T) {
		service := NewService(
			newFakeClaimRepository([]Claim{}, errors.Errorf("")),
			newFakeEvidenceRepositories(PropertyKey{}, Evidence{}, nil),
		)
		_, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Error(t, err)
	})
}
