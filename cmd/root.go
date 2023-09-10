package cmd

import (
	"errors"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:  "command-clearer",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal().Err(err).Msg("Failed to load .env file")
		}

		var discordToken string
		var discordAppId string

		discordTokenVar, _ := cmd.Flags().GetString("discord-token-var")
		if discordTokenVar != "" {
			discordToken = os.Getenv(discordTokenVar)
		} else {
			discordToken, _ = cmd.Flags().GetString("discord-token")
		}

		discordAppIdVar, _ := cmd.Flags().GetString("discord-app-id-var")
		if discordAppIdVar != "" {
			discordAppId = os.Getenv(discordAppIdVar)
		} else {
			discordAppId, _ = cmd.Flags().GetString("discord-app-id")
		}

		if discordToken == "" {
			log.Fatal().Msg("Discord token is required - specify --discord-token or --discord-token-var")
		}

		if discordAppId == "" {
			log.Fatal().Msg("Discord application ID is required - specify --discord-app-id or --discord-app-id-var")
		}

		session, err := discordgo.New("Bot " + discordToken)
		checkError(err, "Failed to create Discord session")

		var deletedCommandCount int

		// TODO: better progress reporting
		globalCommands, err := session.ApplicationCommands(discordAppId, "")
		checkError(err, "Failed to get global commands")

		for _, command := range globalCommands {
			err = session.ApplicationCommandDelete(discordAppId, "", command.ID)
			checkError(err, "Failed to delete command")
			deletedCommandCount++
		}

		guilds, err := session.UserGuilds(100, "", "")
		checkError(err, "Failed to get guilds")

		for _, guild := range guilds {
			guildCommands, err := session.ApplicationCommands(discordAppId, guild.ID)
			checkError(err, "Failed to get guild commands")

			for _, guildCommand := range guildCommands {
				err = session.ApplicationCommandDelete(discordAppId, guildCommand.GuildID, guildCommand.ID)
				checkError(err, "Failed to delete guild command")
				deletedCommandCount++
			}
		}

		log.Info().Msgf("Deleted %d commands successfully!", deletedCommandCount)
	},
}

func init() {
	cobra.OnInitialize(initLogging)

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level")
	rootCmd.Flags().String("discord-token-var", "", "name of the environment variable containing the Discord token")
	rootCmd.Flags().String("discord-token", "", "Discord token")
	rootCmd.Flags().String("discord-app-id-var", "", "name of the environment variable containing the Discord application ID")
	rootCmd.Flags().String("discord-app-id", "", "Discord application ID")
}

func initLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	parsedLogLevel, err := zerolog.ParseLevel(logLevel)
	checkError(err, "Failed to parse log level")
	zerolog.SetGlobalLevel(parsedLogLevel)
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatal().Err(err).Msg(message)
	}
}

func Execute() {
	err := rootCmd.Execute()
	checkError(err, "Failed to run")
}
