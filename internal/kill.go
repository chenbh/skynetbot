package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func (s *state) activate(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	if !s.IsAdmin(m) {
		return command.ErrNotAdmin
	}
	s.active = true
	return nil
}

func (s *state) deactivate(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	if !s.IsAdmin(m) {
		return command.ErrNotAdmin
	}
	s.active = false
	return nil
}
