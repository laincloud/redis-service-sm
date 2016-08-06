package proxy

import (
	"github.com/mijia/sweb/log"
	"syscall"
)

const (
	EPOLLET        = 1 << 31
	MaxEpollEvents = 32
)

type aeApiState struct {
	epfd   int
	skfd   int
	events [MaxEpollEvents]syscall.EpollEvent
	cm     *ConnManager
}

func aeApiStateCreate(cm *ConnManager) *aeApiState {
	var event syscall.EpollEvent
	var events [MaxEpollEvents]syscall.EpollEvent

	fd, err := socket()
	if err != nil {
		log.Error("Create Socker err:", err)
		return nil
	}

	epfd, e := syscall.EpollCreate1(0)
	if e != nil {
		log.Error("epoll_create1: ", events)
		return nil
	}

	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)
	if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
		log.Error("epoll_ctl: ", e)
		return nil
	}
	ae := &aeApiState{epfd: epfd, skfd: fd, events: events, cm: cm}
	return ae

}

func (ae *aeApiState) startAeApiPoll() {
	for {
		nevents, err := syscall.EpollWait(ae.epfd, ae.events[:], -1)
		if err != nil {
			if isEINTR(err) {
				continue
			}
			log.Error("Error creating epoll")
		}
		for ev := 0; ev < nevents; ev++ {
			event := ae.events[ev]
			if (event.Events&syscall.EPOLLERR) != 0 ||
				(event.Events&syscall.EPOLLHUP) != 0 ||
				((event.Events & syscall.EPOLLIN) == 0) {
				/* An error has occured on this fd, or the socket is not
				   ready for reading (why were we notified then?) */
				log.Debug("close:", event.Fd)
				ae.CloseConn(int(event.Fd))
				continue
			} else if int(event.Fd) == ae.skfd {
				ae.Accept()
			} else {
				go ae.handleMessage(int(event.Fd))
			}
		}

	}
}

func (ae *aeApiState) close() {
	syscall.Close(ae.skfd)
	syscall.Close(ae.epfd)
}

func (ae *aeApiState) addEvent(fd int) error {
	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN | EPOLLET
	event.Fd = int32(fd)
	if err := syscall.EpollCtl(ae.epfd, syscall.EPOLL_CTL_ADD, fd, &event); err != nil {
		log.Error("addEvent: ", err)
		syscall.Close(fd)
		return err
	}
	return nil
}

func (ae *aeApiState) delEvent(fd int) error {
	var event syscall.EpollEvent
	event.Fd = int32(fd)
	if err := syscall.EpollCtl(ae.epfd, syscall.EPOLL_CTL_DEL, fd, &event); err != nil {
		log.Error("delEvent: ", err)
		return err
	}
	syscall.Close(fd)
	return nil
}
