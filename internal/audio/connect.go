package audio

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func (s *state) connect(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	channel, err := locateChannel(b.Session(), m)
	if err != nil {
		return err
	}
	vc, err := b.Session().ChannelVoiceJoin(m.GuildID, channel, false, false)
	if err != nil {
		return err
	}

	vc.AddHandler(voiceHandler)

	s.closeChan = make(chan struct{}, 1)
	go record(vc, s.closeChan)

	s.vc = vc
	b.Respond(m, "joining voice channel")
	return nil
}

func locateChannel(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		return "", errors.New("can't locate server")
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			return vs.ChannelID, nil
		}
	}
	return "", errors.New("can't locate voice channel")
}
