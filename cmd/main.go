package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize the logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Get the bot token
	token := systems.BotToken()
	if token == "" {
		log.Fatal().Msg("Bot token is empty")
	}

	// Log the bot token
	log.Info().Str("bot token", token).Send()
}
