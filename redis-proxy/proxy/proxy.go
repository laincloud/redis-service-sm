package proxy

import (
	"github.com/laincloud/redis-libs/network"
	"github.com/laincloud/redis-libs/redislibs"
	"github.com/mijia/sweb/log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Proxy struct {
	cur_client      int32
	server_listener net.Listener
	cond            *sync.Cond
	pool            *network.Pool
}

func NEW() *Proxy {
	var p = new(Proxy)
	p.cur_client = 0
	p.cond = sync.NewCond(&sync.Mutex{})

	co := network.NewConnectOption(ReadTimeoutSec, WriteTimeOutSec, BufferSize)

	p.pool = network.NewConnectionPool(co, func() (net.Conn, error) {
		return net.DialTimeout("tcp", master_addr.String(), time.Second*time.Duration(ConnTimeoutSec))
	}, MaxActive, MaxIdle, PoolIdleTimeOutSec)

	p.pool.SetConnStateTest(func(c *network.Conn) bool {
		if err := c.Write(redislibs.COMMAND_PING); err != nil {
			return false
		}
		// clear test response info
		if _, err := c.ReadAll(); err != nil {
			return false
		}
		return true
	})
	return p
}

func (p *Proxy) StartServer() error {
	var err error
	if err = p.initServer(); err != nil {
		log.Fatal(err.Error())
		return err
	}
	p.initListen()
	return nil
}

func (p *Proxy) initServer() error {
	hostAndPort := "0.0.0.0:" + strconv.Itoa(Port)
	var err error
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	if err != nil {
		return err
	}
	p.server_listener, err = net.ListenTCP("tcp", serverAddr)
	if err != nil {
		return err
	}
	log.Debug("Listening to: ", p.server_listener.Addr().String())
	return nil
}

func (p *Proxy) initListen() {
	for {
		p.cond.L.Lock()
		for p.cur_client >= int32(Max_client) {
			log.Warn("Connection limited")
			p.cond.Wait()
		}
		p.cur_client++
		p.cond.L.Unlock()
		conn, _ := p.server_listener.Accept()
		go p.connectionHandler(conn)
	}
}
