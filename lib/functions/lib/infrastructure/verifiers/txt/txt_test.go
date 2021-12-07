package txt

import (
	"claime-verifier/lib/functions/lib/claim"
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type fakeLookUpper struct {
	recordsLookupper
	fakeLookupTXT func(name string) ([]string, error)
}

func (mock fakeLookUpper) LookupTXT(name string) ([]string, error) {
	return mock.fakeLookupTXT(name)
}

func newFakeLookUpper(records []string, err error) fakeLookUpper {
	return fakeLookUpper{
		fakeLookupTXT: func(name string) ([]string, error) {
			return records, err
		},
	}
}

const (
	addressStr        = "0x81A2863ED122811A1197dB2D9b90a720d73ac81c"
	anotherAddressStr = "0x00142C7D23f0E761E997dsa8eF80244E3D123456"
	mockPropertyID    = "claime.io"
	mockRecord        = recordPrefix + addressStr
	mockAnotherRecord = recordPrefix + anotherAddressStr
)

var (
	address        = common.HexToAddress("0x81A2863ED122811A1197dB2D9b90a720d73ac81c")
	anotherAddress = common.HexToAddress(anotherAddressStr)
)

func TestFind(t *testing.T) {
	t.Run("got evidencial txt records", func(t *testing.T) {
		evidence, err := Client{
			lookupper: newFakeLookUpper([]string{mockRecord, mockAnotherRecord}, nil),
		}.Find(context.Background(), claim.Claim{PropertyID: mockPropertyID})
		assert.Nil(t, err)
		assert.Equal(t, mockPropertyID, evidence.PropertyID)
		assert.Equal(t, address, evidence.EOAs[0])
		assert.Equal(t, mockRecord, evidence.Evidences[0])
		assert.Equal(t, anotherAddress, evidence.EOAs[1])
		assert.Equal(t, mockAnotherRecord, evidence.Evidences[1])
	})
	t.Run("got error if no evidencial txt records found", func(t *testing.T) {
		_, err := Client{
			lookupper: newFakeLookUpper([]string{"txt"}, nil),
		}.Find(context.Background(), claim.Claim{PropertyID: mockPropertyID})
		assert.Equal(t, "no evidencial txt records found", err.Error())
	})
	t.Run("got error if err", func(t *testing.T) {
		mockErrorMsg := "error"
		_, err := Client{
			lookupper: newFakeLookUpper([]string{}, errors.New(mockErrorMsg)),
		}.Find(context.Background(), claim.Claim{PropertyID: mockPropertyID})
		assert.Equal(t, mockErrorMsg, err.Error())
	})
}
