package main

import (
	"github.com/laincloud/redis-service-sm/redis-ctl/monitor"
	"github.com/laincloud/redis-service-sm/redis-ctl/redisctl"
	"github.com/laincloud/redis-service-sm/redis-lainlet/redislainlet"
	"github.com/mijia/sweb/log"
	"os"
)

func main() {
	config_file := "ctl.conf"
	if len(os.Args) == 2 {
		config_file = os.Args[1]
	}
	monitor.ConfigMonitor(config_file)
	redisctl.ConfigCtl(config_file)
	if redisctl.DEBUG {
		log.EnableDebug()
	}
	go redislainlet.StartLainLet()
	go monitor.MonitorServer()
	redisctl.StartRedisMonitor()
}
