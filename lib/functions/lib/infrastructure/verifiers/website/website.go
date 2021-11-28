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

// Find find website ownership evidences
func (c Client) Find(ctx context.Context, cl claim.Claim) (claim.Evidence, error) {
	doc, err := c.scraper.get(cl.PropertyID)
	if err != nil {
		return claim.Evidence{}, err
	}
	evidence := claim.Evidence{
		PropertyID: cl.PropertyID,
	}
	doc.Find("head").Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); name == tagName {
			actualEOAstr, _ := s.Attr("content")
			var actualEvidence string
			if actualEOAstr != "" {
				evidence.EOAs = append(evidence.EOAs, common.HexToAddress(actualEOAstr))
				actualEvidence = fmt.Sprintf(`<meta name="%s" content="%s">`, name, actualEOAstr)
			} else {
				actualEvidence = fmt.Sprintf(`<meta name="%s">`, name)
			}
			evidence.Evidences = append(evidence.Evidences, actualEvidence)
		}
	})
	return evidence, nil
}
