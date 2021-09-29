package slackclient

import (
	"claime-verifier/lib/functions/lib/subscribe"
	"testing"
)

func TestParseFields(t *testing.T) {
	t.Run("parse required fields", func(t *testing.T) {
		t.Run("if no args, all fields are empty", func(t *testing.T) {
			input := &subscribe.Input{}
			parseFields(input, []string{})
		})
		tp := "subscribe"
		network := "rinkeby"
		address := "0xaE2b9f801963891fC1eD72F655De266A7ae34FE8"
		event := "Refund"
		url := "https://raw.githubusercontent.com/bridges-inc/kaleido-core/develop/deployments/rinkeby/AdManager.json"

		cs := []struct {
			name        string
			notExpected func(input *subscribe.Input) bool
			options     []string
		}{
			{
				name: "type", notExpected: func(input *subscribe.Input) bool {
					return input.Type != tp
				}, options: []string{tp},
			},
			{
				name: "network", notExpected: func(input *subscribe.Input) bool {
					return input.Network != network
				}, options: []string{tp, network},
			},
			{
				name: "address", notExpected: func(input *subscribe.Input) bool {
					return input.Address != address
				}, options: []string{tp, network, address},
			},
			{
				name: "event", notExpected: func(input *subscribe.Input) bool {
					return input.Event != event
				}, options: []string{tp, network, address, event},
			},
			{
				name: "abiURL", notExpected: func(input *subscribe.Input) bool {
					return input.AbiURL != url
				}, options: []string{tp, network, address, event, url},
			},
		}
		for _, c := range cs {
			t.Run(c.name+" should be parsed", func(t *testing.T) {
				input := &subscribe.Input{}
				parseFields(input, c.options)
				if c.notExpected(input) {
					t.Error("wrong", c.name)
				}
			})
		}
	})
}
