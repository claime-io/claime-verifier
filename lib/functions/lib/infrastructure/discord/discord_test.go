package discord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

const (
	testGuildID = "guildID"
)

func TestToRegisterContractInput(t *testing.T) {
	in := discordgo.ApplicationCommandInteractionData{
		Options: []*discordgo.ApplicationCommandInteractionDataOption{{Value: "roleId"}, {Value: "contractAddress"}, {Value: "rinkeby"}},
	}
	res := ToRegisterContractInput(in, testGuildID)
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
	t.Run("enable to retrive guildID", func(t *testing.T) {
		if res.GuildID != testGuildID {
			t.Error("wrong guildID")
		}
	})
}
