package bot

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	register("disconnect", disconnect)
}

func disconnect(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	if b.vc != nil {
		b.vc.Disconnect()
		b.vc = nil
	}
	return nil
}
