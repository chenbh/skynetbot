package bot

import (
	"errors"
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(command{
		name: "join",
		args: "",
		help: "join your current voice channel",
		fn:   join,
	})
}

func join(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	channel, err := locateChannel(b.session, m)
	if err != nil {
		return err
	}
	vc, err := b.session.ChannelVoiceJoin(m.GuildID, channel, false, false)
	if err != nil {
		return err
	}

	b.vc = vc
	fmt.Fprintln(out, "joining your voice channel")
	return nil
}

func locateChannel(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		return "", errors.New("can't find which server you're in")
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			return vs.ChannelID, nil
		}
	}
	return "", errors.New("can't find voice channel")
}
