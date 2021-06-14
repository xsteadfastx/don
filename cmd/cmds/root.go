// nolint:gochecknoglobals,gochecknoinits
package cmds

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.xsfx.dev/don"
)

const (
	defaultTimeout = 10 * time.Second
	defaultRetry   = time.Second
)

var (
	timeout time.Duration
	retry   time.Duration
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:  "don [command]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := don.Check(don.Cmd(os.Args[1]), timeout, retry); err != nil {
			log.Fatal().Err(err).Msg("received error")
		}

		log.Info().Msg("ready")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version informations",
	Run: func(cmd *cobra.Command, args []string) {
		// nolint: forbidigo
		fmt.Printf("don %s, commit %s, build on %s\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().DurationVarP(&timeout, "timeout", "t", defaultTimeout, "timeout")
	rootCmd.Flags().DurationVarP(&retry, "retry", "r", defaultRetry, "retry")
}

func Execute() error {
	return rootCmd.Execute()
}
