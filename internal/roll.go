package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register("roll", roll)
}

func roll(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	var err error
	ceil := 100
	if len(args) != 0 {
		ceil, err = strconv.Atoi(args[0])
		if err != nil || ceil <= 0 {
			return errors.New("first argument must be a positive int")
		}
	}

	respond(s, m, fmt.Sprintf("%v", rand.Intn(ceil)))
	return nil
}
