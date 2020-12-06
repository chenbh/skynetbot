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
	RunE:  run,
}

func init() {
	viper.SetEnvPrefix("discord")

	rootCmd.Flags().StringP("token", "t", "", "Discord Bot token (DISCORD_TOKEN)")
	viper.BindPFlag("token", rootCmd.Flags().Lookup("token"))
	viper.BindEnv("token")

	rootCmd.Flags().StringP("admin-role", "", "", "Admin role ID (DISCORD_ADMIN_ROLE)")
	viper.BindPFlag("admin_role", rootCmd.Flags().Lookup("admin-role"))
	viper.BindEnv("admin_role")
}

func run(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")
	if token == "" {
		return errors.New("bot token missing")
	}

	roleID := viper.GetString("admin_role")

	bot, err := bot.NewBot(token, roleID)
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
