package bot

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register("list-audio", listAudio)
}

func listAudio(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	files, err := ioutil.ReadDir("audio")
	if err != nil {
		return err
	}

	message := "available audio:\n"
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".dsa") {
			message += fmt.Sprintf("\t%v\n", strings.TrimSuffix(f.Name(), ".dsa"))
		}
	}
	respond(s, m, message)
	return nil
}
