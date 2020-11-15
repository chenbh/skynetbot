package bot

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(command{
		name: "list-audio",
		args: "",
		help: "list available sound clips",
		fn:   listAudio,
	})
}

func listAudio(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	files, err := ioutil.ReadDir("audio")
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "available audio:")

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".opus") {
			fmt.Fprintf(out, "\t%v\n", strings.TrimSuffix(f.Name(), ".opus"))
		}
	}
	return nil
}
