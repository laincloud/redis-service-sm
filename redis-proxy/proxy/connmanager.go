package proxy

import (
	"github.com/laincloud/redis-libs/network"
)

type ConnManager struct {
	ConnMap map[int]network.IConn
	pool    *network.Pool
}

func NewConnManager(pool *network.Pool) *ConnManager {
	connMap := make(map[int]network.IConn, 0)
	return &ConnManager{ConnMap: connMap, pool: pool}
}

func (cm *ConnManager) NewConn(fd int) {
	if conn, err := cm.pool.FetchConn(); err == nil {
		cm.ConnMap[fd] = conn
	}
}

func (cm *ConnManager) FetchConn(fd int) (network.IConn, error) {
	if _, ok := cm.ConnMap[fd]; ok {
		return cm.ConnMap[fd], nil
	} else {
		if conn, err := cm.pool.FetchConn(); err == nil {
			cm.ConnMap[fd] = conn
			return conn, nil
		} else {
			return nil, err
		}
	}

}

func (cm *ConnManager) CloseConn(fd int) {
	conn, _ := cm.FetchConn(fd)
	delete(cm.ConnMap, fd)
	cm.pool.Finished(conn)
}
