package main

import (
	"bytes"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	endpoint, err := ssm.New().WsEndpoint(context.Background(), "rinkeby")
	if err != nil {
		log.Error("rinkeby not supported", err)
		return
	}
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Error("rinkeby not supported", err)
		return
	}
	contractAbi, err := toAbiFromGithubURL("https://raw.githubusercontent.com/bridges-inc/kaleido-core/develop/deployments/rinkeby/AdManager.json")

	contractAddress := common.HexToAddress("0xaE2b9f801963891fC1eD72F655De266A7ae34FE8")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Error("cant subscribe events", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Error("unexpected error occured: ", err)
		case vLog := <-logs:
			fmt.Println("blockNumber:" + strconv.Itoa(int(vLog.BlockNumber)))
			var mp map[string]interface{} = map[string]interface{}{}
			contractAbi.UnpackIntoMap(mp, "Book", vLog.Data)
			fmt.Println("contractAddress:", vLog.Address.String())
			fmt.Println("eventName:", "Book")
			fmt.Println("txHash:" + vLog.TxHash.String())
			for k, v := range mp {
				fmt.Println("key:" + k + ", value:" + fmt.Sprintf("%v", v))
			}
		}
	}
}

func toAbiFromGithubURL(url string) (abi.ABI, error) {
	abiJSON, err := downloadAbiFromGithub(url)
	if err != nil {
		fmt.Print(err.Error())
		return abi.ABI{}, err
	}
	return toAbiFromJSON(abiJSON)
}

func toAbiFromJSON(js []byte) (abi.ABI, error) {
	maybeAbiRecord, err := abi.JSON(bytes.NewReader(js))
	if err == nil {
		return maybeAbiRecord, nil
	}

	m := make(map[string]interface{})
	if err = json.Unmarshal(js, &m); err != nil {
		fmt.Println("abi from json failed")
		return abi.ABI{}, err
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
		return abi.ABI{}, err
	}
	return fromAbiRecord, err
}

func downloadAbiFromGithub(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error("get abi failed", err)
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
