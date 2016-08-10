package proxy

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/laincloud/redis-libs/network"
	"github.com/mijia/sweb/log"
	"net"
	"syscall"
)

var (
	errRedisDown      = errors.New("-Error redis is down\r\n")
	errRedisDownBytes = []byte(errRedisDown.Error())
)

func (ae *aeApiState) handleMessage(fd int) {
	b, err := network.SyscallRead(fd, BufferSize)
	if err != nil {
		network.SyscallWrite(fd, &errRedisDownBytes, BufferSize)
		return
	}
	if len(bytes.TrimSpace(b)) == 0 {
		return
	}
	if resp, err := ae.Fetcher(fd, &b); err == nil {
		network.SyscallWrite(fd, &resp, BufferSize)
	} else {
		errorsBytes := []byte(err.Error())
		if bytes.HasPrefix(errorsBytes, []byte("-Err")) || bytes.HasPrefix(errorsBytes, []byte("-Error")) {
			network.SyscallWrite(fd, &errorsBytes, BufferSize)
		} else {
			network.SyscallWrite(fd, &errRedisDownBytes, BufferSize)
		}
	}
}

func CounterCmd(b *[]byte) int {
	br := bufio.NewReader(bytes.NewReader(*b))
	rr := network.NewRedisReader(br)
	counts := 0
	for {
		if _, err := rr.ReadObject(); err != nil {
			break
		}
		counts++
	}
	return counts
}

func (ae *aeApiState) Fetcher(fd int, reqs *[]byte) ([]byte, error) {
	redisConn, err := ae.cm.FetchConn(fd)
	if err != nil {
		return nil, err
	}
	if err = redisConn.Write([]byte(*reqs)); err != nil {
		ae.CloseConn(fd)
		return nil, err
	}
	resps := make([]byte, 0)
	counter := CounterCmd(reqs)
	for i := 0; i < counter; i++ {
		if resp, err := redisConn.ReadAll(); err != nil {
			return nil, err
		} else {
			resps = append(resps, resp...)
		}
	}
	return resps, nil
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
