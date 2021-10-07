package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestToInput(t *testing.T) {
	in := discordgo.ApplicationCommandInteractionData{
		Options: []*discordgo.ApplicationCommandInteractionDataOption{{Value: "roleId"}, {Value: "contractAddress"}, {Value: "rinkeby"}},
	}
	res := toInput(in)
	t.Run("enable to retrive roleId", func(t *testing.T) {
		if res.RoleID != "roleId" {
			t.Error("wrong roledid")
		}
	})
	t.Run("enable to retrive contractAddress", func(t *testing.T) {
		if res.ContractAddress != "contractAddress" {
			t.Error("wrong contractAddress")
		}
	})
	t.Run("enable to retrive network", func(t *testing.T) {
		if res.Network != "rinkeby" {
			t.Error("wrong network")
		}
	})
}
