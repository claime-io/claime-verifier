package website

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

type httpScraper struct{}

func (s httpScraper) get(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return &goquery.Document{}, err
	}
	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &goquery.Document{}, err
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(buffer)
	if err != nil {
		return &goquery.Document{}, err
	}

	bufferReader := bytes.NewReader(buffer)
	reader, err := charset.NewReaderLabel(result.Charset, bufferReader)
	if err != nil {
		return &goquery.Document{}, err
	}

	return goquery.NewDocumentFromReader(reader)
}
