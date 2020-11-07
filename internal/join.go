package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register("join", join)
}

func join(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	channel, err := locateChannel(s, m)
	if err != nil {
		return err
	}
	vc, err := s.ChannelVoiceJoin(m.GuildID, channel, false, false)
	if err != nil {
		return err
	}

	b.vc = vc
	respond(s, m, "joining your voice channel")
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
