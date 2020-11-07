package bot

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(command{
		name: "disconnect",
		args: "",
		help: "disconnect from current channel",
		fn:   disconnect,
	})
}

func disconnect(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	if b.vc != nil {
		b.vc.Disconnect()
		b.vc = nil
	}
	return nil
}
