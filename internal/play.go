package bot

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"
)

func init() {
	register("play", play)
}

func play(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	if b.vc == nil {
		return errors.New("need to join voice channel first")
	}

	if len(args) != 1 {
		return errors.New("unexpected amount of arguments")
	}

	filename := args[0]
	buffer, err := load(path.Join("audio", filename+".dsa"))
	if err != nil {
		return err
	}

	if b.stopPlayback != nil {
		b.stopPlayback()
	}

	var ctx context.Context
	ctx, b.stopPlayback = context.WithCancel(b.ctx)

	go transmit(ctx, b, buffer)
	return nil
}

func transmit(ctx context.Context, b *bot, buffer [][]byte) {
	for _, buf := range buffer {
		select {
		case <-ctx.Done():
			return
		default:
			b.vc.OpusSend <- buf
		}
	}
}

func load(filename string) ([][]byte, error) {
	buffer := make([][]byte, 0)

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}

	var opuslen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			break
		}

		if err != nil {
			return nil, fmt.Errorf("reading file: %v", err)
		}

		buf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &buf)
		if err != nil {
			return nil, err
		}

		buffer = append(buffer, buf)
	}

	return buffer, nil
}
