package bot

import (
	"github.com/bwmarrin/discordgo"
)

func setActive(active bool) action {
	return func(b *bot, args []string, m *discordgo.MessageCreate) error {
		err := b.requireAdmin(m)
		if err == nil {
			b.active = active
		}
		return err
	}
}
