package main

import (
	"os"

	"github.com/cirocosta/oklog-docker-plugin/driver"
	"github.com/cirocosta/oklog-docker-plugin/http"
	"github.com/cirocosta/oklog-docker-plugin/oklog"

	"github.com/alexflint/go-arg"
	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/rs/zerolog"
)

type config struct {
	Debug    bool   `arg:"help:whether to activate debug logs"`
	Socket   string `arg:"help:unix socket to listen to"`
	Ingester string `arg:"env:INGESTER,required,help:address of oklog ingester"`
}

var (
	handler = sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	logger  = zerolog.New(os.Stdout).
		With().
		Str("from", "main").
		Logger()
	args = &config{
		Socket:   "oklog",
		Debug:    false,
		Ingester: "",
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
	var (
		err       error
		okLog     oklog.OkLog
		logDriver driver.Driver
	)

	arg.MustParse(args)

	if !args.Debug {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	okLog, err = oklog.New(oklog.Config{
		Host: args.Ingester,
	})
	must(err)

	logDriver, err = driver.New(driver.Config{
		OkLog: &okLog,
	})
	must(err)

	http.Handlers(&handler, &logDriver)

	err = handler.ServeUnix(args.Socket, 0)
	must(err)
}
