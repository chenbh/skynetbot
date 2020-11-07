package bot

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(command{
		name: "roll",
		args: "[CEIL]",
		help: "generate a random number in [0, 100) or [0, CEIL)",
		fn:   roll,
	})
}

func roll(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error {
	var err error
	ceil := 100
	if len(args) != 0 {
		ceil, err = strconv.Atoi(args[0])
		if err != nil || ceil <= 0 {
			return errors.New("first argument must be a positive int")
		}
	}

	fmt.Fprintf(out, "%v", rand.Intn(ceil))
	return nil
}
