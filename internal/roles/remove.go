package roles

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func remove(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	if !b.IsAdmin(m) {
		return command.ErrNotAdmin
	}

	roleID := discordID.FindString(args[0])
	role, err := findRole(b.Session(), m.GuildID, roleID)
	if err != nil {
		return err
	}

	err = b.Session().GuildRoleDelete(m.GuildID, role.ID)
	if err != nil {
		return err
	}

	return b.Respond(m, fmt.Sprintf("removed %v role", role.Mention()))
}
