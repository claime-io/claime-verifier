package registry

import (
	"claime-verifier/lib/functions/lib/claim"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/contracts"
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type (

	// Provider eth provider
	Provider struct {
		cli *ethclient.Client
	}
	// EndpointResolver resolver
	EndpointResolver interface {
		EndpointByNetwork(ctx context.Context, network string) (val string, err error)
	}
)

// NewProvider new Provider
func NewProvider(ctx context.Context, network string, res EndpointResolver) (Provider, error) {
	e, err := res.EndpointByNetwork(ctx, network)
	if err != nil {
		log.Error("resolve endpoint failed", err)
		return Provider{}, err
	}
	client, err := ethclient.Dial(e)
	if err != nil {
		log.Error("provider create failed", err)
		return Provider{}, err
	}
	return Provider{
		cli: client,
	}, nil
}

func callOpts(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{
		Pending: false,
		Context: ctx,
	}
}

func (p Provider) newRegistry() (*contracts.ContractsCaller, error) {
	return contracts.NewContractsCaller(common.HexToAddress(registryAddress()), p.cli)
}

// ClaimsOf claims of eoa
func (p Provider) ClaimsOf(ctx context.Context, eoa common.Address) ([]claim.Claim, error) {
	reg, err := p.newRegistry()
	if err != nil {
		return []claim.Claim{}, err
	}
	keys, _, err := reg.ListClaims(callOpts(ctx), eoa)
	if err != nil {
		return []claim.Claim{}, err
	}
	res := []claim.Claim{}
	for _, key := range keys {
		c, err := reg.AllClaims(callOpts(ctx), key)
		if err != nil {
			return []claim.Claim{}, err
		}
		res = append(res, claim.Claim{
			PropertyType: c.PropertyType,
			PropertyID:   c.PropertyId,
			Evidence:     c.Evidence,
			Method:       c.Method,
		})
	}
	return res, nil
}
