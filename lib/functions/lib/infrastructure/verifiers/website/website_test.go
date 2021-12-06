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
	propertyID        = "https://example.com"
	addressStr        = "0xCdfc500F7f0FCe1278aECb0340b523cD55b3EBbb"
	anotherAddressStr = "0x00142C7D23f0E761E997dsa8eF80244E3D123456"
)

var (
	validEvidence            = fmt.Sprintf(`<meta name="%s" content="%s">`, tagName, addressStr)
	anotherEvidence          = fmt.Sprintf(`<meta name="%s" content="%s">`, tagName, anotherAddressStr)
	multipleEvidence         = anotherEvidence + validEvidence
	address                  = common.HexToAddress(addressStr)
	anotherAddress           = common.HexToAddress(anotherAddressStr)
	emptyContentEvidence     = fmt.Sprintf(`<meta name="%s">`, tagName)
	evidenceHtmlStr          = fmt.Sprintf(`<html><head><meta name="Name" content="Content">%s<title>Title</title></head><body></body></html>`, validEvidence)
	multipleEvidencesHtmlStr = fmt.Sprintf(`<html><head><meta name="Name" content="Content">%s<title>Title</title></head><body></body></html>`, multipleEvidence)
	emptyContentHtmlStr      = fmt.Sprintf(`<html><head><meta name="Name" content="Content">%s<title>Title</title></head><body></body></html>`, emptyContentEvidence)
)

func TestFind(t *testing.T) {
	t.Run("return eoa, actual", func(t *testing.T) {
		evidence, err := Client{
			scraper: newFakeScraper(evidenceHtmlStr, nil),
		}.Find(context.Background(), claim.Claim{PropertyID: propertyID})
		assert.Nil(t, err)
		assert.Equal(t, address, evidence.EOAs[0])
		assert.Equal(t, validEvidence, evidence.Evidences[0])
		assert.Equal(t, propertyID, evidence.PropertyID)
	})
	t.Run("return multiple eoas", func(t *testing.T) {
		evidence, err := Client{
			scraper: newFakeScraper(multipleEvidencesHtmlStr, nil),
		}.Find(context.Background(), claim.Claim{PropertyID: propertyID})
		assert.Nil(t, err)
		assert.Equal(t, anotherAddress, evidence.EOAs[0])
		assert.Equal(t, anotherEvidence, evidence.Evidences[0])
		assert.Equal(t, address, evidence.EOAs[1])
		assert.Equal(t, validEvidence, evidence.Evidences[1])
		assert.Equal(t, propertyID, evidence.PropertyID)
	})
	t.Run("return empty if claim not found", func(t *testing.T) {
		evidence, err := Client{
			scraper: newFakeScraper("", nil),
		}.Find(context.Background(), claim.Claim{PropertyID: propertyID})
		assert.Nil(t, err)
		assert.Empty(t, evidence.EOAs)
		assert.Empty(t, evidence.Evidences)
		assert.Equal(t, propertyID, evidence.PropertyID)
	})
	t.Run("return meta tag with name only if content not found", func(t *testing.T) {
		evidence, err := Client{
			scraper: newFakeScraper(emptyContentHtmlStr, nil),
		}.Find(context.Background(), claim.Claim{PropertyID: propertyID})
		assert.Nil(t, err)
		assert.Empty(t, evidence.EOAs)
		assert.Equal(t, emptyContentEvidence, evidence.Evidences[0])
		assert.Equal(t, propertyID, evidence.PropertyID)
	})
}
