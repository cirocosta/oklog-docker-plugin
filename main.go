package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/logpp/driver"
	"github.com/cirocosta/logpp/http"
	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/rs/zerolog"
)

type config struct {
	Socket string `arg:"help:unix socket to listen to"`
}

var (
	handler = sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	logger  = zerolog.New(os.Stdout).
		With().
		Str("from", "main").
		Logger()
	args = &config{
		Socket: "logpp.socket",
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

	http.Handlers(&h, driver.NewDriver())

	err = handler.ServeUnix("log", 0)
	must(err)
}
