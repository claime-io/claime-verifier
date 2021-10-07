package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestToInput(t *testing.T) {
	in := discordgo.ApplicationCommandInteractionData{
		Options: []*discordgo.ApplicationCommandInteractionDataOption{{Value: "roleId"}, {Value: "contractAddress"}, {Value: 1}},
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
	t.Run("enable to retrive chainId", func(t *testing.T) {
		if res.ChainID != 1 {
			t.Error("wrong chainId")
		}
	})
}
