package proxy

import (
	"errors"
	"github.com/mijia/sweb/log"
	"io"
	"net"
	"time"
)

var (
	errRedisDown = errors.New("-Error redis is down\r\n")
)

func (p *Proxy) connectionHandler(conn net.Conn) {
	log.Debug("receive connection from ", conn.RemoteAddr())
	defer p.disconnect(conn)
	redisConn, err := p.FetchConn()
	if err != nil {
		return
	}
	defer p.disconnect(redisConn)
	isClientClosed := make(chan struct{}, 1)
	isServerClosed := make(chan struct{}, 1)
	go pipe(redisConn, conn, isClientClosed)
	go pipe(conn, redisConn, isServerClosed)
	select {
	case <-isClientClosed:
	case <-isServerClosed:
	}
}

func (p *Proxy) FetchConn() (net.Conn, error) {
	if master_addr == nil {
		return nil, errRedisDown
	}
	return net.DialTimeout("tcp", master_addr.String(),
		time.Second*time.Duration(ConnTimeoutSec))
}

func pipe(dst, src net.Conn, isSrcClosed chan<- struct{}) {
	defer close(isSrcClosed)
	if _, err := io.Copy(dst, src); err != nil {
		log.Debug("copy err:", err)
	}
	isSrcClosed <- struct{}{}
}

func (p *Proxy) disconnect(conn net.Conn) {
	if conn != nil {
		conn.Close()
	}
	p.cond.L.Lock()
	p.cur_client--
	p.cond.L.Unlock()
	p.cond.Signal()
}
