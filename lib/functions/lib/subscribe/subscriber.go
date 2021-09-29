package subscribe

import (
	"claime-verifier/lib/functions/lib/common/log"
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	rinkebyScanBaseURL = "https://rinkeby.etherscan.io/tx/"
	polygonScanBaseURL = "https://polygonscan.com/tx"
)

type (
	// Subscriber subscriber
	Subscriber struct {
		endpoint      wsEndpoint
		repository    Repository
		resolver      AbiResolver
		sender        MsgSender
		ethsubscriber EthSubscriber
	}

	// Event event from smart contract.
	Event struct {
		BlockNumber     uint64
		EventName       string
		ContractAddress string
		TxHash          string
		Link            string
		Parameters      map[string]interface{}
	}

	Network struct {
		Identifier, ScanBaseUrl string
	}

	// EthSubscriber subscriber
	EthSubscriber interface {
		Subscribe(ctx context.Context, subscription Subscription, endpoint string, disp func(ctx context.Context, event Event, sub Subscription) error) error
	}

	// Dispatcher message dispather
	Dispatcher interface {
		Dispatch(ctx context.Context, event Event, sub Subscription) error
	}

	// MsgSender sender
	MsgSender interface {
		Send(ctx context.Context, event Event, subscription Subscription) error
	}

	wsEndpoint struct {
		rinkeby, polygon string
	}

	// Repository repository
	Repository interface {
		Subscribing(ctx context.Context) ([]Subscription, error)
		SubscribingFrom(ctx context.Context, t time.Time) ([]Subscription, error)
	}

	// AbiResolver abi resolver
	AbiResolver interface {
		Resolve(ctx context.Context, url string) (abi.ABI, error)
	}

	// EndpointResolver resolver
	EndpointResolver interface {
		WsEndpoint(ctx context.Context, network string) (val string, err error)
	}

	// Subscription subscription
	Subscription struct {
		Address    string    `dynamo:"PK"`
		EventName  string    `dynamo:"SK"`
		Event      abi.Event `dynamo:"Event"`
		Timestamp  time.Time `dynamo:"Timestamp"`
		Network    string    `dynamo:"Network"`
		WebhookURL string    `dynamo:"WebhookUrl"`
		Abi        abi.ABI   `dynamo:"Abi"`
		ChannelID  string    `dynamo:"ChannelId"`
	}
)

// NewSubscriber new subscriber
func NewSubscriber(ctx context.Context, resolver EndpointResolver, sender MsgSender, repository Repository, ethSubscriber EthSubscriber) (Subscriber, error) {
	rinkeby, err := resolver.WsEndpoint(ctx, NetworkRinkeby)
	if err != nil {
		log.Error("cannot resolve endpoint rinkeby", err)
	}
	polygon, err := resolver.WsEndpoint(ctx, NetworkPolygon)
	if err != nil {
		log.Error("cannot resolve endpoint polygon", err)
	}
	return Subscriber{
		endpoint: wsEndpoint{
			rinkeby: rinkeby,
			polygon: polygon,
		},
		repository:    repository,
		sender:        sender,
		ethsubscriber: ethSubscriber,
	}, err
}

// ScanURLByNetwork scan url
func ScanURLByNetwork(network, address string) string {
	if network == NetworkRinkeby {
		return rinkebyScanBaseURL + address
	}
	if network == NetworkPolygon {
		return polygonScanBaseURL + address
	}
	return ""
}

func newSubscription(input Input, contract abi.ABI) (Subscription, error) {
	ev, err := toEvent(input.Event, contract)
	if err != nil {
		return Subscription{}, err
	}
	return Subscription{
		Address:    input.Address,
		EventName:  input.Event,
		Event:      ev,
		Network:    input.Network,
		Abi:        contract,
		Timestamp:  time.Now(),
		WebhookURL: input.WebhookURL,
		ChannelID:  input.ChannelID,
	}, nil
}

// StartSubscription start subscription
func (s Subscriber) StartSubscription(ctx context.Context) error {
	subs, err := s.loadAllSubscriptions(ctx)
	if err != nil {
		return err
	}
	s.startSubscriptions(ctx, subs)
	return nil
}

func (s Subscriber) startSubscriptions(ctx context.Context, subs []Subscription) {
	var wg sync.WaitGroup
	for i := 0; i < len(subs); i++ {
		wg.Add(1)
		sub := subs[i]
		log.Info("start subscription event " + sub.Event.RawName + " of " + sub.Address)
		go s.startSubscriptionOf(ctx, sub)
	}
	wg.Wait()
}

func (s Subscriber) startSubscriptionsFrom(ctx context.Context, t time.Time) {
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
		}
	}
}

func (s Subscriber) startSubscriptionOf(ctx context.Context, subscription Subscription) error {
	return s.ethsubscriber.Subscribe(ctx, subscription, s.endpoint.rinkeby, func(ctx context.Context, event Event, sub Subscription) error {
		return s.dispatch(ctx, event, sub)
	})
}

func (s Subscriber) dispatch(ctx context.Context, event Event, sub Subscription) error {
	if err := s.sender.Send(ctx, event, sub); err != nil {
		log.Error("msg send failed:", err)
	}
	return nil
}

func newEvent(vLog types.Log, subscription Subscription) Event {
	var mp map[string]interface{} = map[string]interface{}{}
	subscription.Abi.UnpackIntoMap(mp, subscription.Event.RawName, vLog.Data)
	return Event{
		BlockNumber:     vLog.BlockNumber,
		ContractAddress: vLog.Address.String(),
		EventName:       subscription.EventName,
		TxHash:          vLog.TxHash.String(),
		Link:            ScanURLByNetwork(subscription.Network, vLog.TxHash.String()),
		Parameters:      mp,
	}
}

func (s Subscriber) loadAllSubscriptions(ctx context.Context) ([]Subscription, error) {
	return s.repository.Subscribing(ctx)
}
