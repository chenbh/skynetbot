package bot

import (
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

func init() {
	register(command{
		name: "roles",
		args: "",
		help: "list all the roles you are part of",
		fn:   roles,
	})
}

func roles(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	guildRoles, botRole, err := getRoles(b.session, m.GuildID)
	if err != nil {
		return err
	}

	user, err := b.session.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Available roles for %v are:\n", m.Author.Mention())
	fmt.Fprintln(out, "```")

	table := tablewriter.NewWriter(out)
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
	fmt.Fprintln(out, "```")
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
