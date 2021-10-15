package txt

import (
	"claime-verifier/lib/functions/lib/common/log"
	"context"
	"errors"
	"net"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const (
	recordPrefix = "claime-ownership-claim="
)

type (
	// Client client
	Client struct{}
)

// New new client
func New() Client {
	return Client{}
}

// EOA get eoa from domain
func (c Client) EOA(ctx context.Context, domain string) (common.Address, error) {
	txtrecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Error("nslookup failed", err)
	}
	for _, txt := range txtrecords {
		if strings.HasPrefix(txt, recordPrefix) {
			address := strings.ReplaceAll(txt, recordPrefix, "")
			return common.HexToAddress(address), nil
		}
	}
	return common.Address{}, errors.New("no txt records found")
}
