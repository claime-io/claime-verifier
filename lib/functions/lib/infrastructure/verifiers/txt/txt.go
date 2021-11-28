package txt

import (
	"claime-verifier/lib/functions/lib/claim"
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
	Client struct {
		lookupper recordsLookupper
	}
	recordsLookupper interface {
		LookupTXT(name string) ([]string, error)
	}
	dnsService struct{}
)

func (s dnsService) LookupTXT(name string) ([]string, error) {
	return net.LookupTXT(name)
}

// New new client
func New() Client {
	return Client{lookupper: dnsService{}}
}

// Find domain ownership evidences
func (c Client) Find(ctx context.Context, cl claim.Claim) (claim.Evidence, error) {
	txtrecords, err := c.lookupper.LookupTXT(cl.PropertyID)
	if err != nil {
		log.Error("nslookup failed", err)
		return claim.Evidence{}, err
	}
	evidence := claim.Evidence{
		PropertyID: cl.PropertyID,
	}
	for _, txt := range txtrecords {
		if strings.HasPrefix(txt, recordPrefix) {
			actualAddress := strings.ReplaceAll(txt, recordPrefix, "")
			evidence.EOAs = append(evidence.EOAs, common.HexToAddress(actualAddress))
			evidence.Evidences = append(evidence.Evidences, txt)
		}
	}
	if len(evidence.EOAs) == 0 {
		return claim.Evidence{}, errors.New("no evidencial txt records found")
	}
	return evidence, nil
}
