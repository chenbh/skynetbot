package command

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrNotAdmin = errors.New("admin role required")
)

type Bot interface {
	Session() *discordgo.Session
	IsAdmin(*discordgo.MessageCreate) bool
	Respond(*discordgo.MessageCreate, string) error
}

type Action func(b Bot, args []string, m *discordgo.MessageCreate) error

type Cmd struct {
	name string
	args []string
	help string

	fn   Action          // single action
	cmds map[string]*Cmd // action group/family

	cleanup func()
}

func NewAction(name string, args []string, help string, fn Action) *Cmd {
	return &Cmd{
		name: name,
		args: args,
		help: help,
		fn:   fn,
	}
}

func NewGroup(name, help string) *Cmd {
	c := make(map[string]*Cmd)
	return &Cmd{
		name: name,
		help: help,
		cmds: c,
	}
}

func (c *Cmd) AddCommand(cmd *Cmd) {
	if c.cmds == nil {
		panic("adding commands to an action")
	}
	if _, found := c.cmds[cmd.name]; found {
		fmt.Println("duplicate command registered: " + cmd.name)
	}
	c.cmds[cmd.name] = cmd
}

func (c *Cmd) AddCleanup(fn func()) {
	if c.cleanup != nil {
		panic("multiple cleanup added")
	}
	c.cleanup = fn
}

func (c *Cmd) Run(b Bot, args []string, m *discordgo.MessageCreate) error {
	// command is a single action
	if c.fn != nil {
		if len(args) != len(c.args) {
			msg := fmt.Sprintf("%v expected %v args, but got %v", c.name, len(c.args), len(args))
			return b.Respond(m, msg)
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
			return b.Respond(m, err.Error())
		}
		return nil
	} else if args[0] == "help" {
		msg := fmt.Sprintf("available commands for `%v`:\n", c.name)
		msg += c.displayHelp()
		return b.Respond(m, msg)
	} else {
		msg := fmt.Sprintln(out, "unknown command, available commands:")
		msg += c.displayHelp()
		return b.Respond(m, msg)
	}
}

func (c *Cmd) Stop() {
	if c.cleanup != nil {
		c.cleanup()
	}

	if c.cmds != nil {
		for _, v := range c.cmds {
			v.Stop()
		}
	}
}

func (c *Cmd) displayHelp() string {
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
