package audio

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
	"github.com/chenbh/skynetbot/pkg/opusfile"
)

func (s *state) play(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	if s.vc == nil {
		return errors.New("need to join voice channel first")
	}

	if len(args) != 1 {
		return errors.New("unexpected amount of arguments")
	}

	filename := args[0]
	reader, err := load(path.Join("audio", filename+".opus"))
	if err != nil {
		return err
	}

	if s.stopPlayback != nil {
		s.stopPlayback()
	}

	var ctx context.Context
	ctx, s.stopPlayback = context.WithCancel(s.ctx)

	go transmit(ctx, s.vc.OpusSend, reader)
	return nil
}

func transmit(ctx context.Context, out chan<- []byte, reader opusfile.OpusReader) {
	for {
		packet, err := reader.NextPacket()
		if err != nil {
			break
		}

		select {
		case <-ctx.Done():
			return
		default:
			out <- packet
		}
	}
}

func load(filename string) (opusfile.OpusReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}

	reader, err := opusfile.NewOpusReader(file)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
