package audio

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/command"
)

func (s *state) remove(b command.Bot, args []string, m *discordgo.MessageCreate) error {
	filename := args[0]

	err := validateFile(filename)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join("audio", filename+".opus"))
	if err != nil {
		return err
	}

	b.Respond(m, fmt.Sprintf("removed %v", filename))
	return nil
}
