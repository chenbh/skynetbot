package bot

import (
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type action func(b *bot, args []string, m *discordgo.MessageCreate, out io.Writer) error

type command struct {
	name string
	args string
	help string
	fn   action
}

var commands map[string]command

const controlChar = "/"

func register(cmd command) {
	if commands == nil {
		commands = make(map[string]command)
	}

	if _, found := commands[cmd.name]; found {
		panic("duplicate command registered: " + cmd.name)
	}

	commands[cmd.name] = cmd
}

func (b *bot) handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	out := &strings.Builder{}

	if !b.handleAdminCommand(m, out) {
		b.handleCommand(m, out)
	}

	s.ChannelMessageSend(m.ChannelID, out.String())
}

func (b *bot) handleAdminCommand(m *discordgo.MessageCreate, out io.Writer) bool {
	if b.adminChannelID != "" && m.ChannelID == b.adminChannelID && m.Content == "/kill" {
		b.killed = true
	}

	if b.killed {
		fmt.Fprintln(out, "kill switch toggled, no longer handling commands")
		return true
	}

	return false
}

func (b *bot) handleCommand(m *discordgo.MessageCreate, out io.Writer) {
	// log stack trace, similar to http.Serve
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(out, err)
			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("%v\n%s", err, buf)
		}
	}()

	if strings.HasPrefix(m.Content, controlChar) {
		content := strings.TrimLeft(m.Content, controlChar)
		content = strings.Trim(content, " ")

		args := parseArgs(content)

		if len(args) == 0 {
			return
		}

		if cmd, found := commands[args[0]]; found {
			err := cmd.fn(b, args[1:], m, out)
			if err != nil {
				fmt.Fprintf(out, "error: %v", err.Error())
			}
		} else {
			displayHelp(out)
		}
	}
}

func parseArgs(line string) []string {
	// TODO: handle quoted args
	return strings.Fields(line)
}

func displayHelp(out io.Writer) {
	fmt.Fprintln(out, "Unknown command, available commands are:")

	cmds := sort.StringSlice{}
	for _, cmd := range commands {
		cmds = append(cmds, fmt.Sprintf("`%v %v`: %v", cmd.name, cmd.args, cmd.help))
	}

	cmds.Sort()
	for _, v := range cmds {
		fmt.Fprintf(out, "\t%v\n", v)
	}
}
