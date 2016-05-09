package redisctl

import (
	"github.com/laincloud/redis-service-sm/redis-lainlet/redislainlet"
	"github.com/mijia/sweb/log"
	"time"
)

func StartRedisMonitor() {
	for {
		time.Sleep(time.Duration(REDIS_MANAGER_INTERVAL_MS) * time.Millisecond)
		log.Debug("check redis and sentinel status")
		ManagerRedisService(redislainlet.Redis_addrs(), redislainlet.Sentinel_addrs())
	}
}
