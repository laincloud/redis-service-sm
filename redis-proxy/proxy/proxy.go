package proxy

import (
	"github.com/laincloud/redis-libs/network"
	"github.com/laincloud/redis-libs/redislibs"
	"net"
	"time"
)

type Proxy struct {
	pool *network.Pool
	aes  *aeApiState
}

type msgHandler func(msg string) (string, error)

func NewProxy() *Proxy {
	master_addr = redislibs.BuildAddress("127.0.0.1", "6001")

	co := network.NewConnectOption(ReadTimeoutSec, WriteTimeOutSec, BufferSize)
	pool := network.NewConnectionPool(co, func() (net.Conn, error) {
		if master_addr == nil {
			return nil, errRedisDown
		}
		return net.DialTimeout("tcp", master_addr.String(), time.Second*time.Duration(ConnTimeoutSec))
	}, MaxActive, MaxIdle, PoolIdleTimeOutSec)

	pool.SetConnStateTest(func(c *network.Conn) bool {
		if err := c.Write(redislibs.COMMAND_PING); err != nil {
			return false
		}
		// clear test response info
		if _, err := c.ReadAll(); err != nil {
			return false
		}
		return true
	})

	p := &Proxy{pool: pool}

	p.aes = aeApiStateCreate(p.handleMsg)

	return p
}

func (p *Proxy) StartServer() {
	if p.aes == nil {
		return
	}
	p.aes.startAeApiPoll()
}

func (p *Proxy) StopServer() {
	p.aes.close()
}

func (p *Proxy) handleMsg(reqs string) (string, error) {
	redisConn, err := p.pool.FetchConn()
	if err != nil {
		return "", err
	}
	defer p.pool.Finished(redisConn)
	err = redisConn.Write(reqs)
	if err != nil {
		return "", err
	}
	return redisConn.ReadAll()
}
