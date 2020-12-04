package bot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/pkg/opusfile"
)

func audioCommand() *command {
	audio := newGroup("audio", "voice channel oundboard")
	audio.addCommand(newAction(
		"list",
		nil,
		"list available sound clips",
		listAudio,
	))
	audio.addCommand(newAction(
		"upload",
		nil,
		"upload sound clip(s) using attachments, see `audio list`",
		upload,
	))
	audio.addCommand(newAction(
		"remove",
		[]string{"NAME"},
		"remove a sound clip, see `audio list`",
		remove,
	))
	audio.addCommand(newAction(
		"join",
		nil,
		"join you current voice channel",
		join,
	))
	audio.addCommand(newAction(
		"disconnect",
		nil,
		"disconnect from current voice channel",
		disconnect,
	))
	audio.addCommand(newAction(
		"play",
		[]string{"NAME"},
		"plays a sound clip, see `audio list`",
		play,
	))
	audio.addCommand(newAction(
		"stop",
		nil,
		"stop playing the current sound clip",
		stop,
	))
	return audio
}

func listAudio(b *bot, args []string, m *discordgo.MessageCreate) error {
	files, err := ioutil.ReadDir("audio")
	if err != nil {
		return err
	}

	out := &strings.Builder{}
	fmt.Fprintln(out, "available audio:")

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".opus") {
			fmt.Fprintf(out, "\t%v\n", strings.TrimSuffix(f.Name(), ".opus"))
		}
	}
	respond(b.session, m, out.String())
	return nil
}

func upload(b *bot, args []string, m *discordgo.MessageCreate) error {
	for _, a := range m.Attachments {
		err := validateFile(a.Filename)
		if err != nil {
			return err
		}

		// to make it feel more interactive, send messages while processing stuff
		respond(b.session, m, fmt.Sprintf("procccessing %v...", a.Filename))
		path, err := downloadUrl(a.URL, a.Filename)
		if err != nil {
			return fmt.Errorf("downloading: %v", err)
		}

		err = convert(path)
		if err != nil {
			return fmt.Errorf("converting: %v", err)
		}
	}
	// out is only written after the command returns
	respond(b.session, m, "done!")
	return nil
}

func downloadUrl(url, filename string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	path := filepath.Join("audio", filename)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}

func convert(filename string) error {
	ext := filepath.Ext(filename)
	outputFile := strings.TrimRight(filename, ext) + ".opus"

	options := []string{
		"-i", filename,
		"-ar", "48000", // 48khz
		"-ac", "2", // stereo
		"-b:a", "64K", // bitrate 64kbs
		outputFile,
	}
	err := exec.Command("ffmpeg", options...).Run()
	if err != nil {
		return fmt.Errorf("ffmpeg: %v", err)
	}

	return os.Remove(filename)
}

func remove(b *bot, args []string, m *discordgo.MessageCreate) error {
	filename := args[0]
	err := validateFile(filename)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join("audio", filename+".opus"))
	if err != nil {
		return err
	}

	respond(b.session, m, fmt.Sprintf("removed %v", filename))
	return nil
}

func join(b *bot, args []string, m *discordgo.MessageCreate) error {
	channel, err := locateChannel(b.session, m)
	if err != nil {
		return err
	}
	vc, err := b.session.ChannelVoiceJoin(m.GuildID, channel, false, false)
	if err != nil {
		return err
	}

	b.vc = vc
	respond(b.session, m, "joining your voice channel")
	return nil
}

func disconnect(b *bot, args []string, m *discordgo.MessageCreate) error {
	if b.vc != nil {
		b.vc.Disconnect()
		b.vc = nil
	}
	return nil
}

func play(b *bot, args []string, m *discordgo.MessageCreate) error {
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

func stop(b *bot, args []string, m *discordgo.MessageCreate) error {
	if b.vc != nil && b.stopPlayback != nil {
		b.stopPlayback()
		b.stopPlayback = nil
	}
	return nil
}
