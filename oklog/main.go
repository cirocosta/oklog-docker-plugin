package oklog

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"gopkg.in/fatih/pool.v2"
)

type OkLog struct {
	host string
	pool pool.Pool
}

type Config struct {
	Host string
}

func New(cfg Config) (o OkLog, err error) {
	var p pool.Pool

	if cfg.Host == "" {
		err = errors.Errorf("Host must be non-empty")
		return
	}

	p, err = pool.NewChannelPool(2, 15, func() (net.Conn, error) {
		return net.Dial("tcp", cfg.Host)
	})
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create connection pool to %s",
			cfg.Host)
		return
	}

	o.pool = p

	return
}

// TODO what to do when writes start to fail?
func (o *OkLog) Write(line string) (err error) {
	var conn net.Conn

	conn, err = o.pool.Get()
	if err != nil {
		err = errors.Wrapf(err, "failed to get connection from pool")
		return
	}
	defer conn.Close()

	_, err = fmt.Fprintln(conn, line)
	if err != nil {
		err = errors.Wrapf(err, "failed to write line to connection")

		if pc, ok := conn.(*pool.PoolConn); ok {
			pc.MarkUnusable()
			pc.Close()
		}

		return
	}

	return
}

func (o *OkLog) Close() (err error) {
	o.pool.Close()
	return
}
