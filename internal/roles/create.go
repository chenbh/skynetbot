package roles

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func create(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	// err := b.requireAdmin(m)
	// if err != nil {
	// 	return err
	// }

	name := args[0]
	role, err := b.Session().GuildRoleCreate(m.GuildID)
	if err != nil {
		return err
	}

	role, err = b.Session().GuildRoleEdit(
		m.GuildID, role.ID,
		name,
		role.Color, role.Hoist, role.Permissions, role.Mentionable,
	)
	if err != nil {
		return err
	}

	return b.Respond(m, fmt.Sprintf("created %v role", role.Mention()))
}
