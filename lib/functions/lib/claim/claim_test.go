package claim

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
)

type fakeRepository struct {
	Repository
	fakeClaimsOf func(ctx context.Context, eoa common.Address) ([]Claim, error)
}

func (mock fakeRepository) ClaimsOf(ctx context.Context, eoa common.Address) ([]Claim, error) {
	return mock.fakeClaimsOf(ctx, eoa)
}

func newFakeRepository(claims []Claim, err error) fakeRepository {
	return fakeRepository{
		fakeClaimsOf: func(ctx context.Context, eoa common.Address) ([]Claim, error) {
			return claims, err
		},
	}
}

type fakeEvidenceRepository struct {
	EvidenceRepository
	fakeEOA func(ctx context.Context, cl Claim) (EOAOutput, error)
}

func (mock fakeEvidenceRepository) EOA(ctx context.Context, cl Claim) (EOAOutput, error) {
	return mock.fakeEOA(ctx, cl)
}

func newFakeEvidenceRepository(out EOAOutput, err error) fakeEvidenceRepository {
	return fakeEvidenceRepository{
		fakeEOA: func(ctx context.Context, cl Claim) (EOAOutput, error) {
			return out, err
		},
	}
}

func newFakeVerifier(out EOAOutput, err error) fakeEvidenceRepository {
	return fakeEvidenceRepository{
		fakeEOA: func(ctx context.Context, cl Claim) (EOAOutput, error) {
			return out, err
		},
	}
}

func newFakeVerfiers(verifier Verifier, out EOAOutput, err error) map[Verifier]EvidenceRepository {
	return map[Verifier]EvidenceRepository{
		verifier: newFakeVerifier(out, err),
	}
}

var (
	verifyingEOA = common.HexToAddress("0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb")
	mockClaim    = Claim{
		PropertyType: "TestPropertyType",
		PropertyID:   "TestPropertyID",
		Evidence:     "",
		Method:       "TestMethod",
	}
	mockVerifier = Verifier{
		PropertyType: mockClaim.PropertyType,
		Method:       mockClaim.Method,
	}
	mockOutput = EOAOutput{
		Actual:     "TestActual",
		Got:        verifyingEOA,
		PropertyID: mockClaim.PropertyID,
	}
)

func TestVerifyClaims(t *testing.T) {
	t.Run("return verified if matched got EOA", func(t *testing.T) {
		service := NewService(
			newFakeRepository([]Claim{mockClaim}, nil),
			newFakeVerfiers(mockVerifier, mockOutput, nil),
		)
		outputs, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Nil(t, err)
		assert.NotEmpty(t, outputs)
		assert.Equal(t, outputs[0].Result, verified)
	})
	t.Run("return failed if not matched got EOA", func(t *testing.T) {
		service := NewService(
			newFakeRepository([]Claim{mockClaim}, nil),
			newFakeVerfiers(mockVerifier, EOAOutput{
				Actual:     mockOutput.Actual,
				Got:        common.HexToAddress("differentEOA"),
				PropertyID: mockOutput.PropertyID,
			}, nil),
		)
		outputs, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Nil(t, err)
		assert.NotEmpty(t, outputs)
		assert.Equal(t, outputs[0].Result, failed)
	})
	t.Run("return failed if not matched got propertyID", func(t *testing.T) {
		service := NewService(
			newFakeRepository([]Claim{mockClaim}, nil),
			newFakeVerfiers(mockVerifier, EOAOutput{
				Actual:     mockOutput.Actual,
				Got:        mockOutput.Got,
				PropertyID: "differentID",
			}, nil),
		)
		outputs, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Nil(t, err)
		assert.NotEmpty(t, outputs)
		assert.Equal(t, outputs[0].Result, failed)
	})
	t.Run("return failed and message if error", func(t *testing.T) {
		expectedMessage := "Error Message"
		service := NewService(
			newFakeRepository([]Claim{mockClaim}, nil),
			newFakeVerfiers(mockVerifier, EOAOutput{}, errors.Errorf(expectedMessage)),
		)
		outputs, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Nil(t, err)
		assert.NotEmpty(t, outputs)
		assert.Equal(t, outputs[0].Result, failed)
		assert.Equal(t, outputs[0].Actual, expectedMessage)
	})
	t.Run("return unsuported if verifier not found", func(t *testing.T) {
		service := NewService(
			newFakeRepository([]Claim{mockClaim}, nil),
			newFakeVerfiers(Verifier{}, EOAOutput{}, nil),
		)
		outputs, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Nil(t, err)
		assert.NotEmpty(t, outputs)
		assert.Equal(t, outputs[0].Result, unsupported)
	})
	t.Run("return error if failed to get claims", func(t *testing.T) {
		service := NewService(
			newFakeRepository([]Claim{}, errors.Errorf("")),
			newFakeVerfiers(Verifier{}, EOAOutput{}, nil),
		)
		_, err := service.VerifyClaims(context.Background(), verifyingEOA)
		assert.Error(t, err)
	})
}
