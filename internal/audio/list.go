package audio

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func (s *state) list(b command.Bot, args []string, m *discordgo.MessageCreate) error {
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
	b.Respond(m, out.String())
	return nil
}
