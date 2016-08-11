package proxy

import (
	"github.com/laincloud/redis-libs/redislibs"
	"github.com/laincloud/redis-service-sm/redis-lainlet/redislainlet"
	"time"
)

func StartWatcher() {
	go redislainlet.StartLainLet()
	go checkMaster()
}

func checkMaster() {
	for {
		time.Sleep(time.Duration(CHECK_INTERVAL_MS) * time.Millisecond)
		if avail_sentinel_addr != nil {
			if master, err := redislibs.GetMasterAddrByName(avail_sentinel_addr.Host, avail_sentinel_addr.Port, MASTER_NAME_SENTINEL); err == nil {
				master_addr = master
				continue
			}
		}
		avail_sentinel_addr = nil
		for _, s_addr := range redislainlet.Sentinel_addrs() {
			if master, err := redislibs.GetMasterAddrByName(s_addr.Host, s_addr.Port, MASTER_NAME_SENTINEL); err == nil {
				master_addr = master
				avail_sentinel_addr = s_addr
				break
			}
		}
	}
}
