package roles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func join(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	roleID := discordID.FindString(args[0])
	role, err := findRole(b.Session(), m.GuildID, roleID)
	if err != nil {
		return err
	}

	err = b.Session().GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
	if err != nil {
		return err
	}

	return b.Respond(m, "assigned to role")
}
