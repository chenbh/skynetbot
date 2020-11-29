package bot

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/pkg/opusfile"
	"golang.org/x/net/context"
)

func init() {
	register(command{
		name: "play",
		args: "NAME",
		help: "plays a sound clip, see list-audio",
		fn:   play,
	})
}

func play(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	if b.vc == nil {
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

	if b.stopPlayback != nil {
		b.stopPlayback()
	}

	var ctx context.Context
	ctx, b.stopPlayback = context.WithCancel(b.ctx)

	go transmit(ctx, b, reader)
	return nil
}

func transmit(ctx context.Context, b *bot, reader opusfile.OpusReader) {
	for {
		packet, err := reader.NextPacketRaw()
		if err != nil {
			break
		}

		select {
		case <-ctx.Done():
			return
		default:
			b.vc.OpusSend <- packet
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
