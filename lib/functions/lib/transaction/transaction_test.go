package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	rawTx                = "0x02f9018e042f8459682f008459682f0f8344aa2094886a9591d624b5360d22b32c586638734080359380b901640e24c52c000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000000f446973636f7264205573657220494400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012383935323334373635313139313638353232000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012436c61696d6520446973636f7264204170700000000000000000000000000000c0"
	signature            = "0x908854790442a23bbfe42f409725e853ebe7e3e6630bf582a1915c7ad4005ec92bc45ab6642a919eb353df1d7433b2f89a225f21d572f69a00e191bd41bf1a341c"
	messageSignature     = "0x03574b2456ba7e3e4fc9f62004c9160dbebbf7731b9b666372ab0d62dc08e5497ba3c6864b6cdadf8211bac776edeeda3426abbdfa2f01851d79af51a38cfcfd1c"
	message              = "{\"propertyType\":\"Discord User ID\",\"propertyId\":\"895234765119168522\",\"evidence\":\"\",\"method\":\"Claime Discord App\"}"
	expectedAddress      = "0x8dc81F896B38167734ca4ff26b1D20C4c78e9190"
	expectedPropertyType = "Discord User ID"
	expectedPropertyID   = "895234765119168522"
	expectedMethod       = "Claime Discord App"
)

func TestRecoverAddressFromTx(t *testing.T) {
	t.Run("recover address", func(t *testing.T) {
		addr, err := recoverAddressFromTx(rawTx, signature)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectedAddress, addr.Hex())
	})
}

func TestRecoverClaimFromTx(t *testing.T) {
	t.Run("recover claim", func(t *testing.T) {
		res, err := recoverClaimFromTx(rawTx)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectedPropertyType, res.PropertyType)
		assert.Equal(t, expectedPropertyID, res.PropertyId)
		assert.Equal(t, expectedMethod, res.Method)
	})
}

func TestRecoverAddressFromMessage(t *testing.T) {
	t.Run("recover address", func(t *testing.T) {
		addr, err := recoverAddressFromMessage(message, messageSignature)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectedAddress, addr.Hex())
	})
}
func TestRecoverClaimFromMessage(t *testing.T) {
	t.Run("recover address", func(t *testing.T) {
		res, err := recoverClaimFromMessage(message)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectedPropertyType, res.PropertyType)
		assert.Equal(t, expectedPropertyID, res.PropertyId)
		assert.Equal(t, expectedMethod, res.Method)
	})
}
