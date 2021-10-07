package guild

import (
	"github.com/bwmarrin/discordgo"
)

func HasPermissionAdministrator(permission int64) bool {
	return permission&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator
}
