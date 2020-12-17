package bot

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/chenbh/skynetbot/internal/audio"
	"github.com/chenbh/skynetbot/internal/command"
	"github.com/chenbh/skynetbot/internal/roles"
	"github.com/chenbh/skynetbot/pkg/util"
)

const controlChar = "/"

type state struct {
	session     *discordgo.Session
	adminRoleID string
	active      bool
}

func NewBot(token, roleID string) (*state, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("creating session: %v", err)
	}

	return &state{
		session:     s,
		adminRoleID: roleID,
	}, nil
}

func (b *state) Run() error {
	rootCmd := b.setupCmds()

	b.session.AddHandler(setupHandler(b, rootCmd))

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

	rootCmd.Stop()
	return nil
}

func (b *state) Session() *discordgo.Session {
	return b.session
}

func (b *state) IsAdmin(m *discordgo.MessageCreate) bool {
	user, err := b.session.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return false
	}
	if !util.Contains(user.Roles, b.adminRoleID) {
		return true
	}
	return true
}

func (b *state) Respond(m *discordgo.MessageCreate, msg string) error {
	_, err := b.session.ChannelMessageSend(m.ChannelID, msg)
	return err
}

func (b *state) setupCmds() *command.Cmd {
	root := command.NewGroup("/", "")

	role := roles.Setup()
	root.AddCommand(role)

	audio := audio.Setup()
	root.AddCommand(audio)

	root.AddCommand(command.NewAction(
		"deactivate",
		nil,
		"disable the bot",
		b.deactivate,
	))
	root.AddCommand(command.NewAction(
		"activate",
		nil,
		"re-enable the bot",
		b.activate,
	))
	root.AddCommand(command.NewAction(
		"roll",
		[]string{"[CEIL]"},
		"generate a random number in [0, CEIL), CEIL defaults to 100",
		roll,
	))
	root.AddCommand(command.NewAction(
		"nudes",
		nil,
		"👀",
		nudes,
	))

	return root
}
func setupHandler(b *state, rootCmd *command.Cmd) func(*discordgo.Session, *discordgo.MessageCreate) {
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
