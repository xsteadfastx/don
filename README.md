# don

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/go.xsfx.dev/don)
[![Go Report Card](https://goreportcard.com/badge/go.xsfx.dev/don)](https://goreportcard.com/report/go.xsfx.dev/don)

Package don is a little helper if you need to check for the readiness of something.
This could be a command to run (like ssh) or a `db.Ping()` for check of the readiness
of a database container.

![readme](./README.gif)

## Use as commandline tool

Download the tool from the [download page](https://github.com/xsteadfastx/don/releases) or
install via brew:

```go
brew tap xsteadfastx/tap [https://github.com/xsteadfastx/homebrew-tap](https://github.com/xsteadfastx/homebrew-tap)
brew install don
don -t 15m -r 15s -c "ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 root@container"
```

This example checks every 15 seconds if the ssh container is ready. It would timeout with an
error after 15 minutes.

## Use as a library

If you want to use don as a library, it just takes a `func() bool` in `don.Ready()`
have a function that runs the readiness check and returns `true` or `false` if its
ready or not. The second argument is the overall timeout and
the third argument is the check interval. Import it like this:

```go
import go.xsfx.dev/don
```

Doing the readiness check like this:

```go
if err := don.Ready(
	func() bool {
		db, err := sql.Open("oracle", dbConn)
		if err != nil {
			log.Warn().Err(err).Str("dbConn", dbConn).Msg("could not open connection")

			return false
		}

		if err := db.Ping(); err != nil {
			log.Warn().Err(err).Str("dbConn", dbConn).Msg("could not ping")

			return false
		}

		return true
	},
	10*time.Minute, // When to timeout completly.
	30*time.Second, // Whats the time between retries.
	false, // If you want a progressbar.
); err != nil {
	log.Error().Err(err).Msg("received error")
	teardown(pool, resource, tmpState.Name())
	os.Exit(1)
}
```

## Functions

### func [Cmd](/don.go#L72)

`func Cmd(c string) func() bool`

Cmd returns a `func() bool` for working with `don.Ready()`. It executes a command and
returns a true if everything looks fine or a false if there was some kind of error.

### func [Ready](/don.go#L90)

`func Ready(f func() bool, timeout time.Duration, retry time.Duration, bar bool) error`

Ready takes a function that executes something and returns a bool to indicate if
something is ready or not. It returns an error if it timeouts.
