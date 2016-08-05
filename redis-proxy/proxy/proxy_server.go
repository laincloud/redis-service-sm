package proxy

import (
	"errors"
	"github.com/laincloud/redis-libs/network"
	"github.com/laincloud/redis-libs/redislibs"
	"github.com/mijia/sweb/log"
	"net"
	"syscall"
)

var errRedisDown = errors.New("-Error redis is down\r\n")

type msgFetcher func(msg []byte) ([]byte, error)

func (ae *aeApiState) handleMessage(fd int) {
	b, err := network.SyscallRead(fd, BufferSize)
	if err != nil {
		network.SyscallWrite(fd, []byte(err.Error()), BufferSize)
		return
	}
	msg := string(b)
	if msg == redislibs.Pack_command("COMMAND") {
		msg = redislibs.Pack_command("PING")
	}
	if resp, err := ae.fetcher([]byte(msg)); err == nil {
		network.SyscallWrite(fd, resp, BufferSize)
	} else {
		network.SyscallWrite(fd, []byte(err.Error()), BufferSize)
	}
}

func (ae *aeApiState) accept() {
	connFd, _, err := syscall.Accept(ae.skfd)
	if isEINTR(err) {
		ae.accept()
	}
	if err != nil {
		log.Error("accept err: ", err)
		return
	}
	log.Info("new connection:", connFd)
	syscall.SetNonblock(connFd, true)
	ae.addEvent(connFd)
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
