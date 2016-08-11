package network

import (
	"errors"
	"github.com/mijia/sweb/log"
	"io"
	"net"
	"os"
	"syscall"
	"time"
)

var ErrConnNil = errors.New("Connector is Nil")

type IConn interface {
	net.Conn
	ReadAll() ([]byte, error)
	WriteAll(msg []byte) error
	ShouldBeClosed() bool
	StopConn()
}

type Conn struct {
	net.Conn
	cnop     *ConnectOption
	err      error
	stopConn chan struct{}
}

func NewConnect(conn net.Conn, co *ConnectOption) (*Conn, error) {
	if conn == nil {
		return nil, ErrConnNil
	}
	stopConn := make(chan struct{}, 1)
	c := &Conn{Conn: conn, cnop: co, stopConn: stopConn}
	return c, nil
}

func (c *Conn) ShouldBeClosed() bool {
	return c.err != nil
}

func (c *Conn) StopConn() {
	c.stopConn <- struct{}{}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	select {
	case <-c.stopConn:
		err = io.EOF
		log.Error("eof")
		return
	default:
		n, err = c.Conn.Read(b)
		c.err = err
		break
	}
	// ch := make(chan struct{})
	// go func(ch chan struct{}) {
	// 	n, err = c.Conn.Read(b)
	// 	c.err = err
	// 	ch <- struct{}{}
	// }(ch)
	// select {
	// case <-ch:
	// 	break
	// case <-c.stopConn:
	// 	err = io.EOF
	// 	log.Error("eof")
	// 	break
	// }
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	for {
		if n, err = c.Conn.Write(b); err != nil {
			log.Error(err)
			if pe, ok := err.(*os.PathError); ok {
				err = pe.Err
			}
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				continue
			}
		}
		break
	}
	c.err = err
	return
}

func (c *Conn) WriteAll(b []byte) error {
	if c == nil {
		return ErrConnNil
	}
	c.SetWriteDeadline(time.Now().Add(c.cnop.wrteTimeOutSec))
	size := len(b)
	from := 0
	for {
		n, err := c.Conn.Write(b[from:])
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			c.err = err
			break
		}
		from += n
		if from == size {
			break
		}
	}
	return c.err
}

func (c *Conn) ReadAll() ([]byte, error) {
	if c == nil {
		return nil, ErrConnNil
	}
	c.SetReadDeadline(time.Now().Add(c.cnop.readTimeOutSec))
	res := make([]byte, 0, 0)
	bufferSize := c.cnop.bufferSize
	buffer := make([]byte, bufferSize)
	for {
		n, err := c.Conn.Read(buffer)
		if err != nil {
			c.err = err
			return res, err
		}
		res = append(res, buffer[:n]...)
		if n < bufferSize {
			break
		}
	}
	return res, nil
}
