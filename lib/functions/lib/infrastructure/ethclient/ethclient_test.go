package ethclient

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

const (
	testEndpoint = "https://rinkeby.infura.io/v3/d06b171ef3ad461fb7e55d033343eba6"
)

var (
	erc721contract = common.HexToAddress("0xdb2B0073915B8879F06592B03404c1501963e317")
	owner          = common.HexToAddress("0x50414Ac6431279824df9968855181474c919a94B")
)

func TestNewERC721Client(t *testing.T) {
	t.Run("dial succeed", func(t *testing.T) {
		_, err := NewERC721Client(testEndpoint)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("error if with wrong endpoint", func(t *testing.T) {
		_, err := NewERC721Client("https://rinkeby.infura.io/v3/aaa")
		if err == nil {
			t.Error(errors.New("no error occured"))
		}
	})
}

func TestCaller(t *testing.T) {
	t.Run("success with existing contract", func(t *testing.T) {
		cli, _ := NewERC721Client(testEndpoint)
		_, err := cli.Caller(erc721contract)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("fail with non-existing contract", func(t *testing.T) {
		cli, _ := NewERC721Client(testEndpoint)
		notexisting := common.HexToAddress("0xdb2B0073915B8879F06592B03404c1501963e316")
		_, err := cli.Caller(notexisting)
		if err == nil {
			t.Error(errors.New("contract exists"))
		}
	})
	t.Run("fail with erc20 contract", func(t *testing.T) {
		// TODO
		//cli, _ := NewERC721Client(testEndpoint)
		//erc20 := common.HexToAddress("0x2DBBbE49f1e4A7dBb183ebE5d4428715DA720aa0")
		//_, err := cli.Caller(erc20)
		//if err == nil {
		//	t.Error(errors.New("contract exists"))
		//}
	})
}

func TestTokenOwner(t *testing.T) {
	t.Run("true if specified EOA has tokens", func(t *testing.T) {
		cli, _ := NewERC721Client(testEndpoint)
		c, _ := cli.Caller(erc721contract)
		if !c.TokenOwner(owner) {
			t.Error("should be an owner")
		}
	})
	t.Run("true if specified EOA have no tokens", func(t *testing.T) {
		cli, _ := NewERC721Client(testEndpoint)
		c, _ := cli.Caller(erc721contract)
		if c.TokenOwner(common.HexToAddress("0x50414Ac6431279824df9968855181474c919a94C")) {
			t.Error("should not be an owner")
		}
	})
}
