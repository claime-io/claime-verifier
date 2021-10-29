package website

import (
	"claime-verifier/lib/functions/lib/claim"
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/ethereum/go-ethereum/common"
)

const (
	tagName = "claime-ownership-claim"
)

type (
	// Client client
	Client struct {
		scraper scraper
	}
	scraper interface {
		get(url string) (*goquery.Document, error)
	}
)

// New new client
func New() Client {
	return Client{
		scraper: httpScraper{},
	}
}

// EOA get eoa from domain
func (c Client) EOA(ctx context.Context, cl claim.Claim) (claim.EOAOutput, error) {
	doc, err := c.scraper.get(cl.PropertyID)
	if err != nil {
		return claim.EOAOutput{}, err
	}
	var eoa string
	Actual := claim.Actual{
		PropertyID: cl.PropertyID,
	}
	doc.Find("head").Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == tagName {
			eoa, _ = s.Attr("content")
			if eoa != "" {
				Actual.Evidence = fmt.Sprintf(`<meta name="%s" content="%s">`, name, eoa)
			} else {
				Actual.Evidence = fmt.Sprintf(`<meta name="%s">`, name)
			}
		}
	})
	return claim.EOAOutput{
		Actual: Actual,
		EOA:    common.HexToAddress(eoa),
	}, nil
}
