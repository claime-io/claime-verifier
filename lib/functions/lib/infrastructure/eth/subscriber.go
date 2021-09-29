package eth

import (
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/subscribe"
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type (
	Subscriber struct{}
)

func (s Subscriber) Subscribe(ctx context.Context, subscription subscribe.Subscription, endpoint string, disp func(ctx context.Context, event subscribe.Event, sub subscribe.Subscription) error) error {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Error(endpoint+" not supported", err)
		return err
	}
	address := common.HexToAddress(subscription.Address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{subscription.Event.ID}},
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
			if err := disp(ctx, newEvent(vLog, subscription), subscription); err != nil {
				log.Error("msg send failed:", err)
			}
		}
	}
}

func newEvent(vLog types.Log, subscription subscribe.Subscription) subscribe.Event {
	var mp map[string]interface{} = map[string]interface{}{}
	subscription.Abi.UnpackIntoMap(mp, subscription.Event.RawName, vLog.Data)
	return subscribe.Event{
		BlockNumber:     vLog.BlockNumber,
		ContractAddress: vLog.Address.String(),
		EventName:       subscription.EventName,
		TxHash:          vLog.TxHash.String(),
		Link:            subscribe.ScanURLByNetwork(subscription.Network, vLog.TxHash.String()),
		Parameters:      mp,
	}
}
