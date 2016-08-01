package monitor

import (
	"github.com/laincloud/redis-libs/redislibs"
	"github.com/laincloud/redis-service-sm/redis-lainlet/redislainlet"
	"github.com/mijia/sweb/log"
	"strconv"
	"time"
)

func RedisServerMonitor(redis_addrs map[int]*redislibs.Address) {
	for instanceNo, rs_addr := range redis_addrs {
		if r, err := RedisServerMetrics(instanceNo, rs_addr.Host, rs_addr.Port); err == nil {
			log.Debug(r)
			tm := time.Now()
			ReportDatas(r, strconv.FormatInt(tm.Unix(), 10))
		} else {
			log.Error(err.Error())
		}
	}
}

func MonitorServer() {
	for {
		time.Sleep(time.Duration(REDIS_MONITOR_INTERVAL_MS) * time.Millisecond)
		log.Info("monitor info")
		RedisServerMonitor(redislainlet.Redis_addrs())
	}
}
