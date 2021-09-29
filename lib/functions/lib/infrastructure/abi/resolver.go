package abiresolver

import (
	"bytes"
	"claime-verifier/lib/functions/lib/common/log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type (
	// Resolver abi resolver
	Resolver struct {
	}
)

const (
	toAbiFailedErrMsg = "You must specify a valid json url. Note: Tried with a Github URL? Try with a raw one: https://raw.githubusercontent.com/XXX/..."
)

// Resolve resolve url to abi
func (r Resolver) Resolve(ctx context.Context, url string) (abi.ABI, error) {
	if url == "" {
		return abi.ABI{}, nil
	}
	return toAbiFromGithubURL(url)
}

func toAbiFromGithubURL(url string) (abi.ABI, error) {
	abiJSON, err := downloadAbiFromURL(url)
	if err != nil {
		return abi.ABI{}, err
	}
	return toAbiFromJSON(abiJSON)
}

func toAbiFromJSON(js []byte) (abi.ABI, error) {
	maybeAbiRecord, err := abi.JSON(bytes.NewReader(js))
	if err == nil {
		return maybeAbiRecord, nil
	}
	toAbiFailedErr := errors.New(toAbiFailedErrMsg)

	m := make(map[string]interface{})
	if err = json.Unmarshal(js, &m); err != nil {
		log.Error("json unmarshal failed", err)
		return abi.ABI{}, toAbiFailedErr
	}
	if _, ok := m["abi"]; !ok {
		return abi.ABI{}, toAbiFailedErr
	}
	abiRecord := m["abi"]
	b, err := json.Marshal(&abiRecord)
	if err != nil {
		fmt.Println("abi from json failed")
		return abi.ABI{}, err
	}
	fromAbiRecord, err := abi.JSON(bytes.NewReader(b))
	if err != nil {
		fmt.Println("abi from json failed")
		return abi.ABI{}, toAbiFailedErr
	}
	return fromAbiRecord, err
}

func downloadAbiFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error("get abi failed", err)
		return []byte{}, errors.New("specified abiUrl not found")
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
