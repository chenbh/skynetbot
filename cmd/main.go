package main

import (
	"errors"
	"fmt"
	"os"

	bot "github.com/chenbh/skynetbot/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "bot",
	Short: "run the bot",
	Long:  `run a god damned discord bot`,
	RunE:  run,
}

func init() {
	viper.SetEnvPrefix("discord")

	rootCmd.Flags().StringP("token", "t", "", "Discord Bot token (DISCORD_TOKEN)")
	viper.BindPFlag("token", rootCmd.Flags().Lookup("token"))
	viper.BindEnv("token")

	rootCmd.Flags().StringP("admin-channel", "", "", "Admin channel ID (DISCORD_ADMIN_CHANNEL)")
	viper.BindPFlag("admin_channel", rootCmd.Flags().Lookup("admin-channel"))
	viper.BindEnv("admin_channel")
}

func run(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return errors.New("bot token missing")
	}

	channelID := viper.GetString("admin_channel")

	bot, err := bot.NewBot(token, channelID)
	if err != nil {
		return err
	}

	return bot.Run()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
