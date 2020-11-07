package bot

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	register("stop", stop)
}

func stop(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	if b.vc != nil && b.stopPlayback != nil {
		b.stopPlayback()
		b.stopPlayback = nil
	}
	return nil
}
