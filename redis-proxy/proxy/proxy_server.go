package proxy

import (
	"errors"
	"github.com/mijia/sweb/log"
	"net"
)

var errRedisDown = errors.New("-Error redis is down\r\n")

func (p *Proxy) connectionHandler(conn net.Conn) {
	defer p.disconnect(conn)
	if conn == nil {
		return
	}
	connFrom := conn.RemoteAddr().String()
	log.Debug("Connection from: ", connFrom)
	if master_addr == nil {
		talktoClients(conn, errRedisDown.Error())
		return
	}
	redisConn, err := p.pool.FetchConn()
	if err != nil {
		talktoClients(conn, err.Error())
		return
	}
	defer p.pool.Finished(redisConn)
	for {
		reqs, err := readString(conn)
		if err != nil {
			return
		}
		err = redisConn.Write(reqs)
		if err != nil {
			talktoClients(conn, errRedisDown.Error())
			return
		}
		msg, err := redisConn.ReadAll()
		if err != nil {
			talktoClients(conn, errRedisDown.Error())
			return
		}
		talktoClients(conn, msg)
	}
}

func readString(conn net.Conn) (string, error) {
	msg := ""
	ibuf := make([]byte, BufferSize)
	for {
		if length, err := conn.Read(ibuf); err == nil {
			bfstr := string(ibuf[:length])
			msg += bfstr
			if length < BufferSize {
				break
			}
		} else {
			return "", err
		}
	}
	return msg, nil
}
func (p *Proxy) disconnect(conn net.Conn) {
	if conn != nil {
		connFromAddr := conn.RemoteAddr().String()
		conn.Close()
		log.Debug("Closed connection:", connFromAddr)
	}
	p.cond.L.Lock()
	p.cur_client--
	p.cond.L.Unlock()
	p.cond.Signal()
}

func talktoClients(to net.Conn, msg string) {
	_, err := to.Write([]byte(msg))
	if err != nil {
		log.Warn(err.Error())
	}
}
