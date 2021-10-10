package ethclient

import (
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/infrastructure/ethclient/erc721"
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type (
	// ERC721Client client for erc721 contract
	ERC721Client struct {
		client *ethclient.Client
	}

	// Caller caller for erc721caller
	Caller struct {
		*erc721.Erc721Caller
	}
)

// NewERC721Client client for erc721 contract
func NewERC721Client(endpoint string) (ERC721Client, error) {
	client, _ := ethclient.Dial(endpoint)
	_, err := client.NetworkID(context.Background())
	if err != nil {
		log.Error(endpoint+" not supported", err)
		return ERC721Client{}, err
	}
	return ERC721Client{
		client: client,
	}, err
}

// Caller caller for ERC721 contract
func (s ERC721Client) Caller(contract common.Address) (*Caller, error) {
	caller, err := erc721.NewErc721Caller(contract, s.client)
	if err != nil {
		return nil, err
	}
	_, err = caller.BalanceOf(nil, common.HexToAddress("0x61414Ac6431279824df9968855181474c919a94B"))
	if err != nil {
		return nil, err
	}
	return &Caller{caller}, err
}

// TokenOwner is owner of token of specified contract.
func (c *Caller) TokenOwner(owner common.Address) bool {
	balance, err := c.BalanceOf(nil, owner)
	return err == nil && balance.Cmp(common.Big0) > 0
}

// IsOwner given EOA Address has NFT of given contract?
func IsOwner(endpoint string, contractAddress common.Address, eoa common.Address) bool {
	cl, err := NewERC721Client(endpoint)
	if err != nil {
		log.Error("", err)
		return false
	}
	caller, err := cl.Caller(contractAddress)
	if err != nil {
		log.Error("", err)
		return false
	}
	return caller.TokenOwner(eoa)
}
