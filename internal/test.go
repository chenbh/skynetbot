package bot

import (
	"github.com/bwmarrin/discordgo"
)

// experimental shit

func nop(b *bot, args []string, m *discordgo.MessageCreate) error {
	return b.session.MessageReactionAdd(m.ChannelID, m.ID, "👋")
}

func nudes(b *bot, args []string, m *discordgo.MessageCreate) error {
	b.respond(m, "https://github.com/chenbh/skynetbot")
	return b.respond(m, "https://i.redd.it/ihgket4706361.png")
}
