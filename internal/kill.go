package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func setActive(active bool) action {
	return func(b *bot, args []string, m *discordgo.MessageCreate) error {
		if b.isAdmin(m.ChannelID) {
			b.active = active
			respond(b.session, m, fmt.Sprintf("setting active to %v", active))
		} else {
			respond(b.session, m, "you don't have the permissions")
		}
		return nil
	}
}
