package audio

import (
	"context"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

var clips map[int]*clip

type clip struct {
	user      string
	startTime uint64
	endTime   uint64
	packets   []*discordgo.Packet
}

func (c *clip) record(p *discordgo.Packet) {
	if c.startTime == 0 {
		c.startTime = uint64(p.Timestamp)
	}

	c.endTime = uint64(p.Timestamp)
	c.packets = append(c.packets, p)
}

type state struct {
	vc           *discordgo.VoiceConnection
	closeChan    chan struct{}
	stopPlayback context.CancelFunc
	ctx          context.Context
}

func (s *state) cleanup() {
	if s.vc != nil {
		s.vc.Disconnect()
	}
}

func Setup() *command.Cmd {
	clips = make(map[int]*clip, 0)
	s := &state{
		ctx: context.Background(),
	}

	audio := command.NewGroup("audio", "voice channel oundboard")
	audio.AddCleanup(s.cleanup)

	audio.AddCommand(command.NewAction(
		"list",
		nil,
		"list available sound clips",
		s.list,
	))
	audio.AddCommand(command.NewAction(
		"clip",
		nil,
		"clip and upload last 60 seconds of audio in voice channel",
		s.clip,
	))
	audio.AddCommand(command.NewAction(
		"upload",
		nil,
		"upload sound clip(s) using attachments, see `audio list`",
		s.upload,
	))
	audio.AddCommand(command.NewAction(
		"remove",
		[]string{"NAME"},
		"remove a sound clip, see `audio list`",
		s.remove,
	))
	audio.AddCommand(command.NewAction(
		"connect",
		nil,
		"connect to your current voice channel",
		s.connect,
	))
	audio.AddCommand(command.NewAction(
		"disconnect",
		nil,
		"disconnect from current voice channel",
		s.disconnect,
	))
	audio.AddCommand(command.NewAction(
		"play",
		[]string{"NAME"},
		"plays a sound clip, see `audio list`",
		s.play,
	))
	audio.AddCommand(command.NewAction(
		"stop",
		nil,
		"stop playing the current sound clip",
		s.stop,
	))
	return audio
}

func validateFile(path string) error {
	if strings.Contains(path, "/") {
		return errors.New("no special characters allowed")
	}
	return nil
}
