package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

var discordID = regexp.MustCompile("[0-9]+")

func roleCommand() *command {
	role := newGroup("roles", "manipulate server roles")
	role.addCommand(newAction(
		"list",
		nil,
		"list all roles you currently have",
		roleList,
	))
	role.addCommand(newAction(
		"create",
		[]string{"NAME"},
		"create a role, requires admin role",
		roleCreate,
	))
	role.addCommand(newAction(
		"remove",
		[]string{"@ROLE"},
		"remove a role, requires admin role",
		roleRemove,
	))
	role.addCommand(newAction(
		"assign",
		[]string{"@USER", "@ROLE"},
		"assign a user to a role, requires admin role",
		roleAssign,
	))
	role.addCommand(newAction(
		"join",
		[]string{"@ROLE"},
		"join a role",
		roleJoin,
	))
	role.addCommand(newAction(
		"leave",
		[]string{"@ROLE"},
		"leave a role",
		roleLeave,
	))
	return role
}

func roleList(b *bot, args []string, m *discordgo.MessageCreate) error {
	out := &strings.Builder{}
	roles, err := b.roles(m.GuildID)
	if err != nil {
		return err
	}

	user, err := b.session.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "roles for %v:\n", m.Author.String())
	fmt.Fprintln(out, "```")

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"NAME", "@MENTION ID", "MEMBER"})

	for _, role := range roles {
		member := contains(user.Roles, role.ID)
		table.Append([]string{role.Name, role.Mention(), fmt.Sprintf("%v", member)})
	}

	table.Render()
	fmt.Fprintln(out, "```")
	b.respond(m, out.String())
	return nil
}

func roleCreate(b *bot, args []string, m *discordgo.MessageCreate) error {
	err := b.requireAdmin(m)
	if err != nil {
		return err
	}

	name := args[0]
	role, err := b.session.GuildRoleCreate(m.GuildID)
	if err != nil {
		return err
	}

	role, err = b.session.GuildRoleEdit(
		m.GuildID, role.ID,
		name,
		role.Color, role.Hoist, role.Permissions, role.Mentionable,
	)
	if err != nil {
		return err
	}

	return b.respond(m, fmt.Sprintf("created %v role", role.Mention()))
}

func roleRemove(b *bot, args []string, m *discordgo.MessageCreate) error {
	err := b.requireAdmin(m)
	if err != nil {
		return err
	}

	roleID := discordID.FindString(args[0])
	role, err := b.findRole(m.GuildID, roleID)
	if err != nil {
		return err
	}

	err = b.session.GuildRoleDelete(m.GuildID, role.ID)
	if err != nil {
		return err
	}

	return b.respond(m, fmt.Sprintf("removed %v role", role.Mention()))
}

func roleAssign(b *bot, args []string, m *discordgo.MessageCreate) error {
	err := b.requireAdmin(m)
	if err != nil {
		fmt.Println(err)
		return err
	}

	memberID := discordID.FindString(args[0])
	member, err := b.session.GuildMember(m.GuildID, memberID)
	if err != nil {
		return err
	}

	roleID := discordID.FindString(args[1])
	role, err := b.findRole(m.GuildID, roleID)
	if err != nil {
		return err
	}

	err = b.session.GuildMemberRoleAdd(m.GuildID, member.User.ID, role.ID)
	if err != nil {
		return err
	}

	return b.respond(m, "assigned to role")
}

func roleJoin(b *bot, args []string, m *discordgo.MessageCreate) error {
	roleID := discordID.FindString(args[0])
	role, err := b.findRole(m.GuildID, roleID)
	if err != nil {
		return err
	}

	err = b.session.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
	if err != nil {
		return err
	}

	return b.respond(m, "assigned to role")
}

func roleLeave(b *bot, args []string, m *discordgo.MessageCreate) error {
	roleID := discordID.FindString(args[0])
	role, err := b.findRole(m.GuildID, roleID)
	if err != nil {
		return err
	}

	err = b.session.GuildMemberRoleRemove(m.GuildID, m.Author.ID, role.ID)
	if err != nil {
		return err
	}

	return b.respond(m, "assigned to role")
}
