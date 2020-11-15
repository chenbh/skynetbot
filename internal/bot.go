package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	session        *discordgo.Session
	ctx            context.Context
	adminChannelID string

	killed bool

	vc           *discordgo.VoiceConnection
	stopPlayback context.CancelFunc
}

func NewBot(token, channelID string) (*bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("creating session: %v", err)
	}

	return &bot{
		session:        s,
		ctx:            context.Background(),
		adminChannelID: channelID,
	}, nil
}

func (b *bot) Run() error {
	b.session.AddHandler(b.handleCommands)

	b.session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds |
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
