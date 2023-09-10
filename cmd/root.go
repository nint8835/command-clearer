package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:  "command-clearer",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("TODO")
	},
}

func init() {
	cobra.OnInitialize(initLogging)

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level")
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
