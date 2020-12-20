package audio

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func (s *state) disconnect(_ command.Bot, _ []string, _ *discordgo.MessageCreate) error {
	s.cleanup()
	return nil
}
