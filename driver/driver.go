package driver

import (
	"os"

	"github.com/cirocosta/oklog-docker-plugin/docker"
	"github.com/rs/zerolog"
)

type Driver struct {
	logger zerolog.Logger
}

func New() (d Driver) {
	d.logger = zerolog.New(os.Stdout)

	return
}

func (d *Driver) StartLogging(file string, info docker.Info) (err error) {
	d.logger.Info().
		Str("file", file).
		Interface("info", info).
		Msg("start")
	return
}

func (d *Driver) StopLogging(file string) (err error) {
	d.logger.Info().
		Str("file", file).
		Msg("stop")

	return
}
