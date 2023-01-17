package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

const NewLine = 10

type Telnet struct {
	address string
	timeout time.Duration
	conn    net.Conn
	closed  bool
	mux     sync.Mutex
	in      io.ReadCloser
	out     io.Writer
}

func (t *Telnet) setClosed(value bool) {
	t.mux.Lock()
	t.closed = value
	t.mux.Unlock()
}

func (t *Telnet) isClosed() bool {
	return t.closed
}

func (t *Telnet) Close() error {
	return t.conn.Close()
}

func (t *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	t.conn = conn
	return err
}

func (t *Telnet) Send() error {
	scanner := bufio.NewScanner(t.in)
	for {
		if !scanner.Scan() {
			return nil
		}

		if t.isClosed() {
			return errors.New("...Connection was closed by peer")
		}

		bytes := scanner.Bytes()
		if len(bytes) == 0 {
			break
		}

		if _, err := t.conn.Write(append(bytes, NewLine)); err != nil {
			return err
		}
	}

	return nil
}

func (t *Telnet) Receive() error {
	if t.isClosed() {
		return nil
	}

	scanner := bufio.NewScanner(t.conn)
	for {
		if !scanner.Scan() {
			t.setClosed(true)
			return nil
		}

		bytes := scanner.Bytes()
		if len(bytes) == 0 {
			break
		}

		if _, err := t.out.Write(append(bytes, NewLine)); err != nil {
			return err
		}
	}

	return nil
}
