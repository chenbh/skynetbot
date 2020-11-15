package bot

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(command{
		name: "remove",
		args: "NAME",
		help: "remove a sound clip, see list-audio",
		fn:   remove,
	})
}

func remove(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	filename := args[0]
	err := validateFile(filename)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join("audio", filename+".opus"))
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "removed %v", filename)
	return nil
}
