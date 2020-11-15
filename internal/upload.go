package bot

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(command{
		name: "upload",
		args: "",
		help: "upload sound clip(s) using attachments, see list-audio",
		fn:   upload,
	})
}

func upload(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	for _, a := range m.Attachments {
		err := validateFile(a.Filename)
		if err != nil {
			return err
		}

		// to make it feel more interactive, send messages while processing stuff
		b.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("proccessing %v...", a.Filename))
		path, err := download(a.URL, a.Filename)
		if err != nil {
			return fmt.Errorf("downloading: %v", err)
		}

		err = convert(path)
		if err != nil {
			return fmt.Errorf("converting: %v", err)
		}
	}
	// out is only written after the command returns
	fmt.Fprintln(out, "done!")
	return nil
}

func download(url, filename string) (string, error) {
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

func validateFile(path string) error {
	if strings.Contains(path, "/") {
		return errors.New("no special characters allowed")
	}
	return nil
}
