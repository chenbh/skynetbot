package roles

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
	"github.com/chenbh/skynetbot/pkg/util"
	"github.com/olekukonko/tablewriter"
)

func list(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	out := &strings.Builder{}
	roles, err := roles(b.Session(), m.GuildID)
	if err != nil {
		return err
	}

	user, err := b.Session().GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "roles for %v:\n", m.Author.String())
	fmt.Fprintln(out, "```")

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"NAME", "@MENTION ID", "MEMBER"})

	for _, role := range roles {
		member := util.Contains(user.Roles, role.ID)
		table.Append([]string{role.Name, role.Mention(), fmt.Sprintf("%v", member)})
	}

	table.Render()
	fmt.Fprintln(out, "```")
	b.Respond(m, out.String())
	return nil
}
