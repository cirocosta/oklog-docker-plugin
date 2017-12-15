package main

import (
	"os"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/oklog-docker-plugin/driver"
	"github.com/cirocosta/oklog-docker-plugin/http"
	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/rs/zerolog"
)

type config struct {
	Socket string `arg:"help:unix socket to listen to"`
}

var (
	err       error
	handler   = sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	logDriver = driver.New()
	logger    = zerolog.New(os.Stdout).
			With().
			Str("from", "main").
			Logger()
	args = &config{
		Socket: "oklog",
	}
)

func must(err error) {
	if err == nil {
		return
	}

	logger.Error().
		Err(err).
		Msg("main execution failed")

	os.Exit(1)
}

func main() {
	arg.MustParse(args)

	http.Handlers(&handler, &logDriver)

	err = handler.ServeUnix(args.Socket, 0)
	must(err)
}
