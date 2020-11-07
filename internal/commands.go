package bot

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type action func(b *bot, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error

var commands map[string]action

const control_char = "/"

func register(name string, fn action) {
	if commands == nil {
		commands = make(map[string]action)
	}

	if _, found := commands[name]; found {
		panic("duplicate command registered: " + name)
	}

	commands[name] = fn
}

func (b *bot) handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if b.adminChannelID != "" && m.ChannelID == b.adminChannelID && m.Content == "/kill" {
		b.killed = true
	}

	if b.killed {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			respond(s, m, fmt.Sprintln(err))
			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("%v\n%s", err, buf)
		}
	}()

	if strings.HasPrefix(m.Content, control_char) {
		content := strings.Trim(m.Content, control_char)
		args := strings.Split(content, " ")

		if len(args) == 0 {
			return
		}

		if fn, found := commands[args[0]]; found {
			err := fn(b, args[1:], s, m)
			if err != nil {
				respond(s, m, fmt.Sprintf("error: %v", err.Error()))
			}
		} else {
			displayHelp(s, m)
		}
	}
}

func respond(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	s.ChannelMessageSend(m.ChannelID, content)
}

func displayHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := "Unknown command, available commands are: \n"
	for k := range commands {
		message += fmt.Sprintf("\t %v\n", k)
	}
	respond(s, m, message)
}
