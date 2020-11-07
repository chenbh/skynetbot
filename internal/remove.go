package bot

import (
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register("remove", remove)
}

func remove(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	filename := args[0]
	err := validateFile(filename)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join("audio", filename+".dsa"))
	if err != nil {
		return err
	}

	respond(s, m, "audio removed")
	return nil
}
