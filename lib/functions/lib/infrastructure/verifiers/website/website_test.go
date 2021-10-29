package website

import (
	"claime-verifier/lib/functions/lib/claim"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PuerkitoBio/goquery"
	"github.com/ethereum/go-ethereum/common"
)

type fakeScraper struct {
	scraper
	fakeGet func(url string) (*goquery.Document, error)
}

func (mock fakeScraper) get(url string) (*goquery.Document, error) {
	return mock.fakeGet(url)
}

func newFakeScraper(htmlStr string, err error) fakeScraper {
	return fakeScraper{
		fakeGet: func(url string) (*goquery.Document, error) {
			return goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
		},
	}
}

const (
	propertyID      = "https://example.com"
	expectedAddress = "0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"
)

var (
	validEvidence        = fmt.Sprintf(`<meta name="%s" content="%s">`, tagName, expectedAddress)
	emptyContentEvidence = fmt.Sprintf(`<meta name="%s">`, tagName)
	htmlStr              = fmt.Sprintf(`<html><head><meta name="Name" content="Content">%s<title>Title</title></head><body></body></html>`, validEvidence)
	emptyContentHtmlStr  = fmt.Sprintf(`<html><head><meta name="Name" content="Content">%s<title>Title</title></head><body></body></html>`, emptyContentEvidence)
)

func TestEOA(t *testing.T) {
	t.Run("return eoa, actual", func(t *testing.T) {
		got, err := Client{
			scraper: newFakeScraper(htmlStr, nil),
		}.EOA(context.Background(), claim.Claim{PropertyID: propertyID})
		assert.Nil(t, err)
		assert.Equal(t, common.HexToAddress(expectedAddress), got.EOA)
		assert.Equal(t, validEvidence, got.Actual.Evidence)
		assert.Equal(t, propertyID, got.Actual.PropertyID)
	})
	t.Run("return empty if claim not found", func(t *testing.T) {
		got, err := Client{
			scraper: newFakeScraper("", nil),
		}.EOA(context.Background(), claim.Claim{PropertyID: propertyID})
		assert.Nil(t, err)
		assert.Equal(t, common.HexToAddress(""), got.EOA)
		assert.Equal(t, "", got.Actual.Evidence)
		assert.Equal(t, propertyID, got.Actual.PropertyID)
	})
	t.Run("return meta tag with name only if content not found", func(t *testing.T) {
		got, err := Client{
			scraper: newFakeScraper(emptyContentHtmlStr, nil),
		}.EOA(context.Background(), claim.Claim{PropertyID: propertyID})
		assert.Nil(t, err)
		assert.Equal(t, emptyContentEvidence, got.Actual.Evidence)
	})
}
