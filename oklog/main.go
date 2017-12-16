package oklog

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type OkLog struct {
	host string
	conn net.Conn
}

type Config struct {
	Host string
}

func New(cfg Config) (o OkLog, err error) {
	if cfg.Host == "" {
		err = errors.Errorf("Host must be non-empty")
		return
	}

	o.host = cfg.Host

	return
}

func (o *OkLog) Write(line string) (err error) {
	_, err = fmt.Fprintln(o.conn, line)
	if err != nil {
		err = errors.Wrapf(err, "failed to write line to connection")
		return
	}

	return
}

func (o *OkLog) Connect() (err error) {
	var conn net.Conn

	conn, err = net.Dial("tcp", o.host)

	if err != nil {
		err = errors.Wrapf(err, "failed to dial host %s", o.host)
		return
	}

	o.conn = conn

	return
}

func (o *OkLog) Disconnect() (err error) {
	if o.conn != nil {
		err = o.conn.Close()
		return
	}

	return
}
