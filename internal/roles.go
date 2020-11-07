package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

func init() {
	register("roles", roles)
}

const check = "✔"
const cross = "❌"

func roles(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	guildRoles, botRole, err := getRoles(s, m.GuildID)
	if err != nil {
		return err
	}

	user, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}

	message := &strings.Builder{}
	message.WriteString(fmt.Sprintf("Available roles for %v are:\n", m.Author.Mention()))
	message.WriteString("```")

	table := tablewriter.NewWriter(message)
	table.SetHeader([]string{"NAME", "@MENTION ID", "MEMBER"})

	for _, role := range guildRoles {
		// don't display any roles higher than the bot
		if role.Position >= botRole.Position || role.Name == "@everyone" {
			continue
		}

		member := contains(user.Roles, role.ID)
		table.Append([]string{role.Name, role.Mention(), fmt.Sprintf("%v", member)})
	}

	table.Render()
	message.WriteString("```")
	respond(s, m, message.String())
	return nil
}

func getRoles(s *discordgo.Session, gid string) ([]*discordgo.Role, *discordgo.Role, error) {
	roles, err := s.GuildRoles(gid)
	if err != nil {
		return nil, nil, err
	}

	bot, err := s.State.Member(gid, s.State.User.ID)
	if err != nil {
		return nil, nil, err
	}

	var botRole *discordgo.Role
	for _, role := range roles {
		if role.Managed && contains(bot.Roles, role.ID) {
			botRole = role
		}
	}

	return roles, botRole, nil
}
