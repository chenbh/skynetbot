package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

// experimental shit

func nop(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	return b.Session().MessageReactionAdd(m.ChannelID, m.ID, "👋")
}

func nudes(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	b.Respond(m, "https://github.com/chenbh/skynetbot")
	return b.Respond(m, "https://i.redd.it/ihgket4706361.png")
}
