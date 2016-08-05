package proxy

import (
	"github.com/laincloud/redis-libs/network"
	"github.com/laincloud/redis-libs/redislibs"
	"github.com/mijia/sweb/log"
	"net"
	"time"
)

type Proxy struct {
	pool *network.Pool
	aes  *aeApiState
}

func NewProxy() *Proxy {
	master_addr = redislibs.BuildAddress("127.0.0.1", "6001")

	co := network.NewConnectOption(ReadTimeoutSec, WriteTimeOutSec, BufferSize)
	pool := network.NewConnectionPool(co, func() (network.IConn, error) {
		if master_addr == nil {
			return nil, errRedisDown
		}
		if conn, err := net.DialTimeout("tcp", master_addr.String(),
			time.Second*time.Duration(ConnTimeoutSec)); err == nil {
			return network.NewRedisConn(conn, co)
		} else {
			return nil, err
		}
	}, MaxActive, MaxIdle, PoolIdleTimeOutSec)

	pool.SetConnStateTest(func(c network.IConn) bool {
		if err := c.Write([]byte(redislibs.COMMAND_PING)); err != nil {
			return false
		}
		// clear test response info
		if _, err := c.ReadAll(); err != nil {
			return false
		}
		return true
	})

	p := &Proxy{pool: pool}

	p.aes = aeApiStateCreate(p.redisMsgFetcher)

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

func (p *Proxy) redisMsgFetcher(reqs []byte) ([]byte, error) {
	redisConn, err := p.pool.FetchConn()
	if err != nil {
		return nil, err
	}
	defer p.pool.Finished(redisConn)
	if err = redisConn.Write([]byte(reqs)); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return redisConn.ReadAll()
}
