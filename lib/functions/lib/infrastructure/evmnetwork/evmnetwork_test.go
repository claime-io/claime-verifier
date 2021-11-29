package evmnetwork

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	t.Run("Equals to string", func(t *testing.T) {
		assert.True(t, Rinkeby.Equals("rinkeby"))
		assert.True(t, Mumbai.Equals("mumbai"))
		assert.True(t, Mainnet.Equals("mainnet"))
		assert.True(t, Polygon.Equals("polygon"))
	})
}
