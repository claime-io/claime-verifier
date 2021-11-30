package registry

import (
	"claime-verifier/lib/functions/lib/claim"
	"claime-verifier/lib/functions/lib/infrastructure/evmnetwork"
	"context"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

type (
	Repository struct {
		providers []Provider
	}
)

var (
	supportedChainsDev  = []evmnetwork.Network{evmnetwork.Rinkeby, evmnetwork.Mumbai}
	supportedChainsProd = []evmnetwork.Network{evmnetwork.Mainnet, evmnetwork.Polygon}
)

func NewRepository(ctx context.Context, network string, res EndpointResolver) (Repository, error) {
	if network != "" {
		provider, err := NewProvider(ctx, network, res)
		if err != nil {
			return Repository{}, err
		}
		return Repository{
			providers: []Provider{provider},
		}, err
	}

	env := os.Getenv("EnvironmentId")
	var chains []evmnetwork.Network
	if env == "prod" {
		chains = supportedChainsProd
	} else {
		chains = supportedChainsDev
	}
	var providers []Provider
	for _, chain := range chains {
		provider, err := NewProvider(ctx, chain.ToString(), res)
		if err != nil {
			return Repository{}, err
		}
		providers = append(providers, provider)
	}
	return Repository{
		providers: providers,
	}, nil
}

func (r Repository) ClaimsOf(ctx context.Context, eoa common.Address) ([]claim.Claim, error) {
	var claims []claim.Claim
	for _, provider := range r.providers {
		res, err := provider.ClaimsOf(ctx, eoa)
		if err != nil {
			return claims, err
		}
		claims = append(claims, res...)
	}
	return claims, nil
}
