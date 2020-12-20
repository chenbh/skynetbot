package audio

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func (s *state) connect(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	if s.vc != nil {
		err := s.disconnect(b, args, m)
		if err != nil {
			return err
		}
	}

	channel, err := locateChannel(b.Session(), m)
	if err != nil {
		return err
	}
	vc, err := b.Session().ChannelVoiceJoin(m.GuildID, channel, false, false)
	if err != nil {
		return err
	}

	s.vc = vc
	s.setupRecording()

	b.Respond(m, "joining voice channel")
	return nil
}

func (s *state) setupRecording() {
	s.vc.AddHandler(s.voiceHandler())

	s.doneRecording = make(chan struct{}, 0)
	s.recordings = make(map[string]*clip)
	s.ssrc = make(map[int]string)

	go s.receiveAudio()
	go s.gcAudio()
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
