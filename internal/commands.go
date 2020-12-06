package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type action func(b *bot, args []string, m *discordgo.MessageCreate) error

type command struct {
	name string
	args []string
	help string

	fn   action
	cmds map[string]*command
}

func newAction(name string, args []string, help string, fn action) *command {
	return &command{
		name: name,
		args: args,
		help: help,
		fn:   fn,
	}
}

func newGroup(name, help string) *command {
	c := make(map[string]*command)
	return &command{
		name: name,
		help: help,
		cmds: c,
	}
}

func (c *command) addCommand(cmd *command) {
	if c.cmds == nil {
		panic("adding commands to an action")
	}
	if _, found := c.cmds[cmd.name]; found {
		fmt.Println("duplicate command registered: " + cmd.name)
	}
	c.cmds[cmd.name] = cmd
}

func (c *command) Run(b *bot, args []string, m *discordgo.MessageCreate) error {
	// command is a single action
	if c.fn != nil {
		if len(args) != len(c.args) {
			msg := fmt.Sprintf("%v expected %v args, but got %v", c.name, len(c.args), len(args))
			return b.respond(m, msg)
		}

		return c.fn(b, args, m)
	}

	// command is a group of commands
	out := &strings.Builder{}
	if cmd, found := c.cmds[args[0]]; found {
		var newArgs []string
		if len(args) > 1 {
			newArgs = args[1:]
		}

		err := cmd.Run(b, newArgs, m)
		if err != nil {
			return b.respond(m, err.Error())
		}
		return nil
	} else if args[0] == "help" {
		msg := fmt.Sprintf("available commands for `%v`:\n", c.name)
		msg += c.displayHelp()
		return b.respond(m, msg)
	} else {
		msg := fmt.Sprintln(out, "unknown command, available commands:")
		msg += c.displayHelp()
		return b.respond(m, msg)
	}
}

func (c *command) displayHelp() string {
	out := &strings.Builder{}
	cmds := sort.StringSlice{}
	for _, cmd := range c.cmds {
		if cmd.cmds != nil {
			nestedHelp := strings.ReplaceAll(cmd.displayHelp(), "\t", "\t\t")
			nestedHelp = strings.TrimRight(nestedHelp, "\n")
			cmds = append(cmds, fmt.Sprintf("`%v`:\n%v", cmd.name, nestedHelp))
		} else {
			usage := strings.Join(append([]string{cmd.name}, cmd.args...), " ")
			cmds = append(cmds, fmt.Sprintf("`%v`: %v", usage, cmd.help))
		}
	}

	cmds.Sort()
	for _, v := range cmds {
		fmt.Fprintf(out, "\t%v\n", v)
	}
	return out.String()
}
