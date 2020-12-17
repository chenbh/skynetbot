package audio

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func (s *state) stop(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	if s.vc != nil && s.stopPlayback != nil {
		s.stopPlayback()
		s.stopPlayback = nil
	}
	return nil
}
