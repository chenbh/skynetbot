package audio

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func (s *state) disconnect(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	if s.vc != nil {
		s.closeChan <- struct{}{}
		s.vc.Disconnect()
		s.vc = nil
	}
	return nil
}
