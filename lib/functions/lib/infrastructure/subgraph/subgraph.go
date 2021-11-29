package subgraph

import (
	"claime-verifier/lib/functions/lib/claim"
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	graphql "github.com/hasura/go-graphql-client"
)

type (
	service struct {
		client subgraphClient
	}
	subgraphClient interface {
		claims(ctx context.Context, eoa common.Address) (listClaims, error)
	}
	graphqlClient struct {
		client *graphql.Client
	}
	listClaims struct {
		Claims []claimData `graphql:"claims(first: 1000, where: { claimer: $claimer, removed: false })"`
	}
	claimData struct {
		ID           graphql.ID
		Claimer      graphql.String
		PropertyType graphql.String
		PropertyID   graphql.String
		Method       graphql.String
		Evidence     graphql.String
		Network      graphql.String
		Removed      graphql.Boolean
	}
)

// New subgraph service
func New(url string) service {
	return newFromClient(graphqlClient{
		client: graphql.NewClient(url, nil),
	},
	)
}

func newFromClient(client subgraphClient) service {
	return service{
		client: client,
	}
}

func (svc service) ClaimsOf(ctx context.Context, eoa common.Address) ([]claim.Claim, error) {
	claims, err := svc.client.claims(ctx, eoa)
	if err != nil {
		return []claim.Claim{}, err
	}
	var res []claim.Claim
	for _, each := range claims.Claims {
		res = append(res, claim.Claim{
			PropertyType: string(each.PropertyType),
			PropertyID:   string(each.PropertyID),
			Method:       string(each.Method),
			Evidence:     string(each.Evidence),
			Network:      string(each.Network),
		})
	}
	return res, nil
}

func (c graphqlClient) claims(ctx context.Context, eoa common.Address) (listClaims, error) {
	variables := map[string]interface{}{
		"claimer": graphql.String(strings.ToLower(eoa.Hex())),
	}
	var res listClaims
	err := c.client.Query(ctx, &res, variables)
	if err != nil {
		return listClaims{}, err
	}
	return res, nil
}
