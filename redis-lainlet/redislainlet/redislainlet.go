package redislainlet

import (
	api "github.com/laincloud/lainlet/api/v2"
	"github.com/laincloud/lainlet/client"
	"github.com/laincloud/redis-libs/redislibs"
	"golang.org/x/net/context"
	"strconv"
	"sync"
	"time"
)

var (
	redis_addrs    = make(map[int]*redislibs.Address)
	sentinel_addrs = make(map[int]*redislibs.Address)
	rwlock         sync.RWMutex
)

func StartLainLet() {
	info := new(api.GeneralPodGroup)
	// get request
	c := client.New(LAINLET_ADDR)
	// watch request
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30) // 30 seconds timeout
	for {
		time.Sleep(time.Duration(ERROR_IDLE_TIME) * time.Millisecond)
		ch, err := c.Watch(WATHCER, ctx)
		if err != nil {
			continue
		}

		for {
			select {
			case event, ok := <-ch:
				if !ok {
					break
				}
				if event.Id != 0 { // id == 0 means error-event or heartbeat
					if err := info.Decode(event.Data); err == nil {
						handleLainLetInfo(info)
					}
				}
			}

		}
	}

}

func handleLainLetInfo(gpd *api.GeneralPodGroup) {
	rwlock.Lock()
	defer rwlock.Unlock()
	datas := gpd.Data
	for _, pg := range datas {
		for _, pd := range pg.Pods {
			if pd.ProcName == PROC_REDIS_NAME {
				redis_addrs[pd.InstanceNo] = redislibs.BuildAddress(pd.IP, strconv.Itoa(pd.Port))
			} else if pd.ProcName == PROC_SENTILE_NAME {
				sentinel_addrs[pd.InstanceNo] = redislibs.BuildAddress(pd.IP, strconv.Itoa(pd.Port))
			}
		}
	}
}

func Redis_addrs() map[int]*redislibs.Address {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return redis_addrs
}

func Sentinel_addrs() map[int]*redislibs.Address {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return sentinel_addrs
}
