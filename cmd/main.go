package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.xsfx.dev/don/cmd/cmds"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := cmds.Execute(); err != nil {
		log.Fatal().Err(err).Msg("received error")
	}
}
