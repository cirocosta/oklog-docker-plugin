package driver

import (
	"context"
	"encoding/binary"
	"io"
	"os"
	"path"
	"sync"
	"syscall"

	"github.com/cirocosta/oklog-docker-plugin/docker"
	"github.com/cirocosta/oklog-docker-plugin/oklog"

	"github.com/docker/docker/api/types/plugins/logdriver"
	protoio "github.com/gogo/protobuf/io"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tonistiigi/fifo"
)

type Driver struct {
	logger zerolog.Logger

	oklog *oklog.OkLog
	logs  map[string]*logPair
	mu    sync.Mutex
}

type Config struct {
	OkLog *oklog.OkLog
}

type logPair struct {
	stream io.ReadCloser
	info   docker.Info
	active bool
}

func New(cfg Config) (d Driver, err error) {
	if cfg.OkLog == nil {
		err = errors.Errorf("OkLog must be non-nil")
		return
	}

	d.logger = zerolog.New(os.Stdout)
	d.logs = make(map[string]*logPair)
	d.oklog = cfg.OkLog

	return
}

func (d *Driver) StartLogging(file string, info docker.Info) (err error) {
	var name = path.Base(file)

	d.mu.Lock()
	_, exists := d.logs[name]
	if exists {
		d.mu.Unlock()
		err = errors.Errorf(
			"logger for %q already exists", file)
		return
	}
	d.mu.Unlock()

	d.logger.Info().
		Str("file", file).
		Interface("info", info).
		Msg("start")

	stream, err := fifo.OpenFifo(
		context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		err = errors.Wrapf(err,
			"error opening logger fifo: %q",
			file)
		return
	}

	lp := &logPair{
		info:   info,
		stream: stream,
		active: true,
	}
	d.logs[name] = lp

	go d.ConsumeLog(lp)

	return
}

func (d *Driver) ConsumeLog(lp *logPair) {
	var (
		buf logdriver.LogEntry
		dec = protoio.NewUint32DelimitedReader(
			lp.stream, binary.BigEndian, 1e6)
		err error
	)

	defer dec.Close()
	defer d.ShutdownLogPair(lp)

	for {
		if !lp.active {
			d.logger.Debug().
				Str("id", lp.info.ContainerID).
				Msg("shutting down logger goroutine due to stop request")
			return
		}

		err = dec.ReadMsg(&buf)
		if err != nil {
			if err == io.EOF {
				d.logger.Debug().
					Err(err).
					Str("id", lp.info.ContainerID).
					Msg("shutting down logger goroutine due to file EOF")
				return
			}

			d.logger.Warn().
				Err(err).
				Str("id", lp.info.ContainerID).
				Msg("error reading from FIFO, trying to continue")

			dec = protoio.NewUint32DelimitedReader(
				lp.stream, binary.BigEndian, 1e6)
			continue
		}

		err = d.DoSomethingWithLog(lp, buf.Line)
		if err != nil {
			d.logger.Warn().
				Str("id", lp.info.ContainerID).
				Err(err).
				Msg("error logging message, dropping it and continuing")
		}

		buf.Reset()
	}
}

func (d *Driver) DoSomethingWithLog(lp *logPair, line []byte) (err error) {
	err = d.oklog.Write(string(line))
	if err != nil {
		err = errors.Wrapf(err, "failed to write line")
		return
	}

	return
}

func (d *Driver) ShutdownLogPair(lp *logPair) {
	if lp.stream != nil {
		lp.stream.Close()
	}

	lp.active = false
}

func (d *Driver) StopLogging(file string) (err error) {
	d.logger.Info().
		Str("file", file).
		Msg("stop")

	var name = path.Base(file)

	d.mu.Lock()
	defer d.mu.Unlock()

	lp, ok := d.logs[name]

	if !ok {
		d.logger.Warn().
			Str("file", file).
			Msg("failed to stop logging - file not active")
		return
	}

	lp.active = false
	delete(d.logs, path.Base(file))

	return
}
