package network

import (
	"bufio"
	"github.com/mijia/sweb/log"
	// "io"
	"os"
	"strconv"
	"syscall"
	"time"
)

const (
	SIMPLE_STRING = '+'
	BULK_STRING   = '$'
	INTEGER       = ':'
	ARRAY         = '*'
	ERROR         = '-'

	SYM_CRLF = "\r\n"
)

func SyscallWrite(fd int, b *[]byte, aeBufferSize int) error {
	from := 0
	to := 0
	respSize := len(*b)
	for {
		if from+aeBufferSize < respSize {
			to = from + aeBufferSize
		} else {
			to = respSize
		}
		if n, err := syscall.Write(fd, (*b)[from:to]); err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			return err
		} else {
			from += n
		}
		if from == respSize {
			break
		}
	}
	return nil
}

type NetFd struct {
	fd int
}

func NewNetFd(fd int) *NetFd {
	return &NetFd{fd: fd}
}

func (f *NetFd) Read(b []byte) (n int, err error) {
	// n, err = syscall.Read(f.fd, b)
	for {
		n, err = syscall.Read(f.fd, b)
		if err != nil {
			log.Info("err:", err)
			if err == syscall.EAGAIN {
				time.Sleep(10 * time.Microsecond)
				continue
			}
		}
		break
	}
	if _, ok := err.(syscall.Errno); ok {
		err = os.NewSyscallError("read", err)
	}
	return
}

func (f *NetFd) Write(b []byte) (n int, err error) {
	size := len(b)
	for {
		n, err = syscall.Write(f.fd, b)
		if err != nil {
			log.Info("err:", err)
			if err == syscall.EAGAIN {
				time.Sleep(10 * time.Millisecond)
				continue
			}
		} else if n < size {
			log.Info(n, " : ", size)
			size -= n
			continue
		}

		break
	}
	// log.Info(size, "  write:", n)
	return
}

func SyscallRead(fd int, aeBufferSize int) ([]byte, error) {
	// fr := NewNetFd(fd)
	// br := bufio.NewReader(fr)
	// redisReader := NewRedisReader(br)
	msg := make([]byte, 0)
	var err error
	bf := make([]byte, aeBufferSize)
	for {
		if n, e := syscall.Read(fd, bf); e == nil {
			if n <= 0 {
				break
			}
			msg = append(msg, bf[:n]...)
		} else {
			log.Error(e)
			if e == syscall.EAGAIN || e == syscall.EWOULDBLOCK {
				break
			}
			err = e
		}
		// if bf, e := redisReader.ReadObject(); e == nil {
		// 	msg = append(msg, bf...)
		// } else {
		// 	if pe, ok := e.(*os.PathError); ok {
		// 		e = pe.Err
		// 	}
		// 	if e == syscall.EAGAIN || e == syscall.EWOULDBLOCK || e == io.EOF {
		// 		break
		// 	}
		// 	err = e
		// 	break
		// }
	}
	return msg, err
}

type RedisReader struct {
	br *bufio.Reader
}

func NewRedisReader(br *bufio.Reader) *RedisReader {
	return &RedisReader{br: br}
}

func (r *RedisReader) ReadObject() ([]byte, error) {
	res := make([]byte, 0)
	err := r.readObject(&res)
	return res, err
}

func (r *RedisReader) readObject(res *[]byte) error {
	line, err := r.readLine()
	if err != nil {
		return err
	}
	*res = append(*res, []byte(line)...)
	switch line[0] {
	case SIMPLE_STRING, INTEGER, ERROR:
		return nil
	case BULK_STRING:
		return r.readBulkString(&line, res)
	case ARRAY:
		return r.readArray(&line, res)
	default:
		return BAD_ELEMENT_ERR
	}
}

func (r *RedisReader) readLine() ([]byte, error) {
	p, err := r.br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *RedisReader) getCount(line *[]byte) (int, error) {
	if !((*line)[len(*line)-1] == '\n' && (*line)[len(*line)-2] == '\r') {
		return 0, BAD_ELEMENT_ERR
	}
	return strconv.Atoi(string((*line)[1 : len(*line)-2]))
}

func (r *RedisReader) readBulkString(line *[]byte, res *[]byte) error {
	size, err := r.getCount(line)
	if err != nil {
		return err
	}
	if size == -1 {
		return nil
	}
	if size < 0 {
		return BAD_ELEMENT_ERR
	}
	size = size + 2
	buffer, err := r.br.Peek(size)
	r.br.Discard(size)
	*res = append(*res, buffer...)
	return err
}

func (r *RedisReader) readArray(line *[]byte, res *[]byte) error {
	// Get number of array elements.
	count, err := r.getCount(line)
	if err != nil {
		return err
	}
	// Read `count` number of RESP objects in the array.
	for i := 0; i < count; i++ {
		err := r.readObject(res)
		if err != nil {
			return err
		}
	}
	return nil
}
