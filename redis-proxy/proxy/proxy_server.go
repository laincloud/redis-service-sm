package proxy

import (
	"bytes"
	"errors"
	"github.com/laincloud/redis-libs/network"
	"github.com/laincloud/redis-libs/redislibs"
	"github.com/mijia/sweb/log"
	"net"
	"strings"
	"syscall"
)

var errRedisDown = errors.New("-Error redis is down\r\n")

func (ae *aeApiState) handleMessage(fd int) {
	b, err := network.SyscallRead(fd, BufferSize)
	if err != nil {
		network.SyscallWrite(fd, []byte(errRedisDown.Error()), BufferSize)
		return
	}
	msg := string(b)
	if strings.Trim(msg, " ") == "" {
		return
	}
	if msg == redislibs.Pack_command("COMMAND") {
		msg = redislibs.Pack_command("PING")
	}
	counter := bytes.Count(b, []byte(redislibs.SYM_STAR))
	for i := 0; i < counter; i++ {
		if resp, err := ae.Fetcher(fd, []byte(msg)); err == nil {
			network.SyscallWrite(fd, resp, BufferSize)
		} else {
			network.SyscallWrite(fd, []byte(errRedisDown.Error()), BufferSize)
		}
	}

}

func (ae *aeApiState) Fetcher(fd int, reqs []byte) ([]byte, error) {
	redisConn, err := ae.cm.FetchConn(fd)
	if err != nil {
		return nil, err
	}
	if err = redisConn.Write([]byte(reqs)); err != nil {
		log.Error(err.Error())
		ae.CloseConn(fd)
		return nil, err
	}
	return redisConn.ReadAll()
}

func (ae *aeApiState) Accept() {
	fd, _, err := syscall.Accept(ae.skfd)
	if isEINTR(err) {
		ae.Accept()
	}
	if err != nil {
		log.Fatal("accept err: ", err)
		return
	}
	log.Debug("new connection:", fd)
	syscall.SetNonblock(fd, true)
	ae.addEvent(fd)
	ae.cm.NewConn(fd)
}

func (ae *aeApiState) CloseConn(fd int) {
	ae.cm.CloseConn(fd)
	ae.delEvent(fd)
}

func socket() (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}
	if err = syscall.SetNonblock(fd, true); err != nil {
		return -1, err
	}
	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	addr := syscall.SockaddrInet4{Port: Port}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	syscall.Bind(fd, &addr)
	syscall.Listen(fd, 32)

	return fd, nil
}

func isEINTR(err error) bool {
	if err == nil {
		return false
	}
	errno, ok := err.(syscall.Errno)
	if ok && errno == syscall.EINTR {
		return true
	}
	return false
}
