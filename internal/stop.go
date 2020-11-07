package bot

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(command{
		name: "stop",
		args: "",
		help: "stop playing the current sound clip",
		fn:   stop,
	})
}

func stop(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	if b.vc != nil && b.stopPlayback != nil {
		b.stopPlayback()
		b.stopPlayback = nil
	}
	return nil
}
