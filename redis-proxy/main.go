package main

import (
	"github.com/laincloud/redis-service-sm/redis-proxy/proxy"
	"github.com/mijia/sweb/log"
	"os"
	// "time"
)

func main() {
	if len(os.Args) == 2 {
		proxy.Load_config(os.Args[1])
	} else {
		proxy.Load_config("proxy.conf")
	}
	if proxy.DEBUG {
		log.EnableDebug()
	}
	proxy.StartWatcher()
	p := proxy.NewProxy()
	p.StartServer()

}
