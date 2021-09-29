package subscribe

import (
	abiresolver "claime-verifier/lib/functions/lib/infrastructure/abi"
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// Subscribe subscribe
	Subscribe = "subscribe"
	// Unsubscribe unsubscribe
	Unsubscribe = "unsubscribe"
	// NetworkRinkeby rinkeby
	NetworkRinkeby = "rinkeby"
	// NetworkPolygon polygon
	NetworkPolygon = "polygon"
)

type (
	// Registrar service
	Registrar struct {
		messenger  Messenger
		repository RegistrarDataSource
	}

	// RegistrarDataSource data source
	RegistrarDataSource interface {
		Subscribe(ctx context.Context, sub Subscription) error
		UnSubscribe(ctx context.Context, in Input) error
	}

	// Messenger messenger application
	Messenger interface {
		ToInputFromBody(body string) (Input, error)
		ErrorOutput(err error) Output
		ToOutput(in Input) Output
	}

	// Input input
	Input struct {
		Type       string
		Address    string
		Event      string
		Network    string
		AbiURL     string
		WebhookURL string
		ChannelID  string
	}

	// Output output
	Output interface {
		Parse() ([]byte, error)
	}
)

var (
	// AllowedTypes allowed Types for subscription.
	AllowedTypes = []string{Subscribe, Unsubscribe}
	// SupportedNetworks supported networks
	SupportedNetworks = []string{NetworkRinkeby, NetworkPolygon}
)

// NewRegistrar new registrar
func NewRegistrar(messenger Messenger, repository RegistrarDataSource) Registrar {
	return Registrar{
		messenger:  messenger,
		repository: repository,
	}
}

// validType is type valid
func validType(val string) bool {
	return valid(val, AllowedTypes)
}

func valid(val string, supported []string) bool {
	for _, s := range supported {
		if val == s {
			return true
		}
	}
	return false
}

// validNetwork is network supported
func validNetwork(val string) bool {
	return valid(val, SupportedNetworks)
}

// Register register a new  subscription
func (s Registrar) Register(ctx context.Context, requestbody string) (Output, error) {
	input, err := s.messenger.ToInputFromBody(requestbody)
	if err != nil {
		return s.messenger.ErrorOutput(err), err
	}
	if err := input.validate(); err != nil {
		return s.messenger.ErrorOutput(err), nil
	}
	if input.Type == Subscribe {
		return s.registerSubscription(ctx, input)
	}
	return s.messenger.ErrorOutput(nil), err
}

func (s Registrar) registerSubscription(ctx context.Context, input Input) (Output, error) {
	abi, err := abiresolver.Resolver{}.Resolve(ctx, input.AbiURL)
	if err != nil {
		return s.messenger.ErrorOutput(err), err
	}
	sub, err := newSubscription(input, abi)
	if err != nil {
		return s.messenger.ErrorOutput(err), err
	}
	if err := s.repository.Subscribe(ctx, sub); err != nil {
		return s.messenger.ErrorOutput(err), err
	}
	return s.messenger.ToOutput(input), nil
}

func toEvent(eventName string, contractAbi abi.ABI) (abi.Event, error) {
	evs := []string{}
	for _, v := range contractAbi.Events {
		evs = append(evs, v.Name)
		if strings.Contains(v.Name, eventName) {
			return v, nil
		}
	}
	errstr := "Event with name " + eventName + " not found. AvailableEvents: "
	for _, e := range evs {
		errstr += "\n-"
		errstr += e
	}
	return abi.Event{}, errors.New(errstr)
}

func (s Registrar) registerUnsubscription(ctx context.Context, input Input) (string, error) {
	return "", nil
}

func (input Input) validate() error {
	if !validType(input.Type) {
		return errors.New("Subscription type is invalid. Allowed values:" + strings.Join(AllowedTypes, ", ") + ". Got:" + input.Type)
	}
	if !validNetwork(input.Network) {
		return errors.New("Not supported network. Allowed values:" + strings.Join(SupportedNetworks, ", ") + ". Got:" + input.Network)
	}
	if input.Address == "" {
		return missingRequiredFieldError("Address")
	}
	if !common.IsHexAddress(input.Address) {
		return errors.New("Contract address " + input.Address + " is not hex address. Example:" + "0xaE2b9f801963891fC1eD72F655De266A7ae34FE8.")
	}
	return nil
}

func missingRequiredFieldError(field string) error {
	return errors.New("Missing required field:" + field)
}
