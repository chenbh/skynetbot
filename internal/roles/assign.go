package roles

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func assign(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	if !b.IsAdmin(m) {
		return command.ErrNotAdmin
	}

	memberID := discordID.FindString(args[0])
	member, err := b.Session().GuildMember(m.GuildID, memberID)
	if err != nil {
		return err
	}

	roleID := discordID.FindString(args[1])
	role, err := findRole(b.Session(), m.GuildID, roleID)
	if err != nil {
		return err
	}

	err = b.Session().GuildMemberRoleAdd(m.GuildID, member.User.ID, role.ID)
	if err != nil {
		return err
	}

	return b.Respond(m, "assigned to role")
}
