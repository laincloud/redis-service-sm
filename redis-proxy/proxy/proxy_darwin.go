package proxy

import (
	"github.com/mijia/sweb/log"
	"syscall"
	"time"
)

const (
	InitKqueueEvents = 32
)

type aeApiState struct {
	kq      int
	skfd    int
	events  []syscall.Kevent_t
	fetcher msgFetcher
	timeout *syscall.Timespec
}

func aeApiStateCreate(fetcher msgFetcher) *aeApiState {

	fd, err := socket()
	if err != nil {
		log.Error("Create Socker err:", err)
		return nil
	}
	kq, err := syscall.Kqueue()
	if err != nil {
		log.Error("Error creating Kqueue descriptor!")
		return nil
	}
	// configure timeout
	events := make([]syscall.Kevent_t, InitKqueueEvents)
	timeout := syscall.Timespec{
		Sec:  0,
		Nsec: 0,
	}
	ae := &aeApiState{skfd: fd, kq: kq, events: events, fetcher: fetcher, timeout: &timeout}
	ae.addEvent(fd)
	return ae
}

func (ae *aeApiState) startAeApiPoll() {
	for {
		nev, err := syscall.Kevent(ae.kq, nil, ae.events, ae.timeout)
		if err != nil {
			if isEINTR(err) {
				continue
			}
			log.Info("Error creating kevent")
		}
		if nev == 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		for i := 0; i < nev; i++ {
			if ae.events[i].Ident == uint64(ae.skfd) {
				ae.accept()
			} else if ae.events[i].Flags&(syscall.EV_EOF|syscall.EV_ERROR) > 0 {
				fd := int(ae.events[i].Ident)
				log.Debug("close:", fd)
				ae.delEvent(fd)
			} else {
				ae.handleMessage(int(ae.events[i].Ident))
			}
		}
	}
}

func (ae *aeApiState) close() {
	if ae == nil {
		return
	}
	syscall.Close(ae.skfd)
	syscall.Close(ae.kq)
}

func (ae *aeApiState) addEvent(fd int) error {
	ev := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}
	if _, err := syscall.Kevent(ae.kq, []syscall.Kevent_t{ev}, nil, ae.timeout); err != nil {
		log.Error("addEvent err:", err)
		syscall.Close(fd)
		return err
	}
	return nil
}

func (ae *aeApiState) delEvent(fd int) error {
	ev := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_DELETE,
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}
	if _, err := syscall.Kevent(ae.kq, []syscall.Kevent_t{ev}, nil, ae.timeout); err != nil {
		log.Error("delEvent err:", err)
		return err
	}
	syscall.Close(fd)
	return nil
}
