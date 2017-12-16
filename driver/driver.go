package driver

import (
	"context"
	"io"
	"os"
	"path"
	"sync"
	"syscall"

	"github.com/cirocosta/oklog-docker-plugin/docker"
	"github.com/tonistiigi/fifo"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Driver struct {
	logger zerolog.Logger

	logs map[string]*logPair
	mu   sync.Mutex
}

type logPair struct {
	stream io.ReadCloser
	info   docker.Info
}

func New() (d Driver) {
	d.logger = zerolog.New(os.Stdout)
	return
}

func (d *Driver) StartLogging(file string, info docker.Info) (err error) {
	var name = path.Base(file)

	d.mu.Lock()
	_, exists := d.logs[name]
	if exists {
		d.mu.Unlock()
		err = errors.Errorf("logger for %q already exists", file)
		return
	}
	d.mu.Unlock()

	d.logger.Info().
		Str("file", file).
		Interface("info", info).
		Msg("start")

	stream, err := fifo.OpenFifo(context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		return errors.Wrapf(err, "error opening logger fifo: %q", file)
	}

	d.logs[name]  = &logPair{
		info: info,
		stream: stream,
	}

	return
}

func (d *Driver) StopLogging(file string) (err error) {
	d.logger.Info().
		Str("file", file).
		Msg("stop")

	return
}
