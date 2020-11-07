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
	register("upload", upload)
}

func upload(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	for _, a := range m.Attachments {
		err := validateFile(a.Filename)
		if err != nil {
			return err
		}

		respond(s, m, fmt.Sprintf("downloading %v...", a.Filename))
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

	respond(s, m, "proccessing files...")
	err := exec.Command("./convert.sh").Run()
	if err != nil {
		return err
	}

	respond(s, m, "done!")
	return nil
}

func validateFile(path string) error {
	if strings.Contains(path, "/") {
		return errors.New("🖕 no special characters allowed")
	}
	return nil
}
