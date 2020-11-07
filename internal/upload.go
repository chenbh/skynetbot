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
		b.session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("downloading %v...", a.Filename))
		resp, err := http.Get(a.URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err := os.Create(filepath.Join("audio", a.Filename))
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}

	b.session.ChannelMessageSend(m.ChannelID, "proccessing files...")
	err := exec.Command("./convert.sh").Run()
	if err != nil {
		return err
	}

	// out is only written after the command returns
	fmt.Fprintln(out, "done!")
	return nil
}

func validateFile(path string) error {
	if strings.Contains(path, "/") {
		return errors.New("no special characters allowed")
	}
	return nil
}
