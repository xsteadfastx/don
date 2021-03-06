// Package don is a little helper if you need to check for the readiness of something.
// This could be a command to run (like ssh) or a `db.Ping()` for check of the readiness
// of a database container.
//
// (image/readme) ./README.gif
//
// Use as commandline tool
//
// Download the tool from the (download page) https://github.com/xsteadfastx/don/releases or
// install via brew:
//
// 	brew tap xsteadfastx/tap https://github.com/xsteadfastx/homebrew-tap
// 	brew install don
// 	don -t 15m -r 15s -c "ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 root@container"
//
// This example checks every 15 seconds if the ssh container is ready. It would timeout with an
// error after 15 minutes.
//
// Use as a library
//
// If you want to use don as a library, it just takes a `func() bool` in `don.Ready()`
// have a function that runs the readiness check and returns `true` or `false` if its
// ready or not. The second argument is the overall timeout and
// the third argument is the check interval. Import it like this:
//
//	import go.xsfx.dev/don
//
// Doing the readiness check like this:
//
//	if err := don.Ready(
//		func() bool {
//			db, err := sql.Open("oracle", dbConn)
//			if err != nil {
//				log.Warn().Err(err).Str("dbConn", dbConn).Msg("could not open connection")
//
//				return false
//			}
//
//			if err := db.Ping(); err != nil {
//				log.Warn().Err(err).Str("dbConn", dbConn).Msg("could not ping")
//
//				return false
//			}
//
//			return true
//		},
//		10*time.Minute, // When to timeout completly.
//		30*time.Second, // Whats the time between retries.
//		false, // If you want a progressbar.
//	); err != nil {
//		log.Error().Err(err).Msg("received error")
//		teardown(pool, resource, tmpState.Name())
//		os.Exit(1)
//	}
//
package don

import (
	"errors"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

var errTimeout = errors.New("timeout")

// Cmd returns a `func() bool` for working with `don.Ready()`. It executes a command and
// returns a true if everything looks fine or a false if there was some kind of error.
func Cmd(c string) func() bool {
	return func() bool {
		cmd := exec.Command("sh", "-c", c)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Warn().Err(err).Msg("cmd has error")

			return false
		}

		return true
	}
}

// Ready takes a function that executes something and returns a bool to indicate if
// something is ready or not. It returns an error if it timeouts.
func Ready(f func() bool, timeout time.Duration, retry time.Duration, bar bool) error {
	chReady := make(chan struct{})

	go func() {
		for {
			if f() {
				chReady <- struct{}{}

				return
			}

			if bar {
				d := int64(retry / time.Second)
				bar := progressbar.Default(d)

				for i := int64(0); i < d; i++ {
					if err := bar.Add(1); err != nil {
						log.Error().Err(err).Msg("could not add to bar")
					}

					time.Sleep(time.Second)
				}
			} else {
				<-time.After(retry)
				log.Info().Msg("retrying")
			}
		}
	}()

	select {
	case <-chReady:
		return nil

	case <-time.After(timeout):
		return errTimeout
	}
}
