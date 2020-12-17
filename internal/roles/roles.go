package roles

import (
	"errors"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
	"github.com/chenbh/skynetbot/pkg/util"
)

var discordID = regexp.MustCompile("[0-9]+")

func Setup() *command.Cmd {
	role := command.NewGroup("roles", "manipulate server roles")
	role.AddCommand(command.NewAction(
		"list",
		nil,
		"list all roles you currently have",
		list,
	))
	role.AddCommand(command.NewAction(
		"create",
		[]string{"NAME"},
		"create a role",
		create,
	))
	role.AddCommand(command.NewAction(
		"remove",
		[]string{"@ROLE"},
		"remove a role, requires admin role",
		remove,
	))
	role.AddCommand(command.NewAction(
		"assign",
		[]string{"@USER", "@ROLE"},
		"assign a user to a role, requires admin role",
		assign,
	))
	role.AddCommand(command.NewAction(
		"join",
		[]string{"@ROLE"},
		"join a role",
		join,
	))
	role.AddCommand(command.NewAction(
		"leave",
		[]string{"@ROLE"},
		"leave a role",
		leave,
	))
	return role
}

func roles(s *discordgo.Session, gid string) ([]*discordgo.Role, error) {
	allRoles, err := s.GuildRoles(gid)
	if err != nil {
		return nil, err
	}

	bot, err := s.State.Member(gid, s.State.User.ID)
	if err != nil {
		return nil, err
	}

	var botRole *discordgo.Role
	for _, role := range allRoles {
		if role.Managed && util.Contains(bot.Roles, role.ID) {
			botRole = role
		}
	}

	roles := make([]*discordgo.Role, 0)
	for _, role := range allRoles {
		if role.Position < botRole.Position && role.Name != "@everyone" {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

func findRole(s *discordgo.Session, gid, rid string) (*discordgo.Role, error) {
	roles, err := roles(s, gid)
	if err != nil {
		return nil, err
	}

	var role *discordgo.Role
	for _, r := range roles {
		if r.ID == rid {
			role = r
		}
	}
	if role == nil {
		return nil, errors.New("cannot manage role")
	}

	return role, nil
}
