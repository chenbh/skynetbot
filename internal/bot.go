package bot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	session     *discordgo.Session
	ctx         context.Context
	adminRoleID string

	active bool

	vc           *discordgo.VoiceConnection
	stopPlayback context.CancelFunc
}

const controlChar = "/"

func NewBot(token, roleID string) (*bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("creating session: %v", err)
	}

	return &bot{
		session:     s,
		ctx:         context.Background(),
		adminRoleID: roleID,
	}, nil
}

func (b *bot) Run() error {
	b.session.AddHandler(setupHandler(b))

	b.session.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsGuilds |
			discordgo.IntentsGuildMessages |
			discordgo.IntentsGuildVoiceStates |
			discordgo.IntentsGuildMessageReactions,
	)

	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("opening session: %v", err)
	}
	defer b.session.Close()

	fmt.Println("running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if b.vc != nil {
		return b.vc.Disconnect()
	}
	return nil
}

func (b *bot) requireAdmin(m *discordgo.MessageCreate) error {
	user, err := b.session.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}
	if !contains(user.Roles, b.adminRoleID) {
		return errors.New("need admin role to run this command")
	}
	return nil
}

func (b *bot) respond(m *discordgo.MessageCreate, msg string) error {
	_, err := b.session.ChannelMessageSend(m.ChannelID, msg)
	return err
}

func setupHandler(b *bot) func(*discordgo.Session, *discordgo.MessageCreate) {
	rootCmd := setupCmds()

	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// log stack trace, similar to http.Serve
		defer func() {
			if err := recover(); err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprint(err))
				buf := make([]byte, 8196)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Printf("%v\n%s", err, buf)
			}
		}()

		if !strings.HasPrefix(m.Content, controlChar) {
			return
		}

		content := strings.TrimLeft(m.Content, controlChar)
		content = strings.Trim(content, " ")

		args := parseArgs(content)

		if len(args) == 0 {
			return
		}

		rootCmd.Run(b, args, m)
	}
}

func parseArgs(line string) []string {
	// TODO: handle quoted args
	return strings.Fields(line)
}

func setupCmds() *command {
	root := newGroup("", "")

	role := roleCommand()
	root.addCommand(role)

	audio := audioCommand()
	root.addCommand(audio)

	root.addCommand(newAction(
		"kill",
		nil,
		"disable the bot",
		setActive(false),
	))
	root.addCommand(newAction(
		"roll",
		[]string{"[CEIL]"},
		"generate a random number in [0, 100) or [0, CEIL)",
		roll,
	))

	return root
}
