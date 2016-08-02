package proxy

import (
	"fmt"
	"github.com/laincloud/redis-libs/network"
	"github.com/laincloud/redis-libs/redislibs"
	// "github.com/mijia/sweb/log"
	"net"
	"testing"
	"time"
)

// func Test_LainLet(t *testing.T) {
// 	// LainLet()
// 	fmt.Println(runtime.GOOS)
// }

// func Test_Kqueue(t *testing.T) {
// 	defer log.Info("out")
// 	Load_config("../proxy.conf")
// 	if DEBUG {
// 		log.EnableDebug()
// 	}
// 	for {
// 		p := NewProxy()
// 		p.StartServer()
// 		p.StopServer()
// 		time.Sleep(60 * time.Second)
// 	}
// }

func Test_Connection(t *testing.T) {
	master_addr = redislibs.BuildAddress("127.0.0.1", "8889")

	c, _ := net.DialTimeout("tcp", master_addr.String(), time.Second*time.Duration(ConnTimeoutSec))
	rc, err := network.NewRedisConn(c, network.NewConnectOption(5, 5, 1024))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cmd := redislibs.Pack_command("keys", "*")
	counter := 0
	for i := 0; i < 10000; i++ {
		rc.Write([]byte(cmd))
		if res, err := rc.ReadAll(); err == nil {
			if len(res) < 10 {
				fmt.Println("i:", i)
				fmt.Println("res:", string(res[:]))
				counter++
			}
		} else {
			fmt.Println("i:", i)
			fmt.Println("err:", err.Error())
			counter++
		}
	}
	fmt.Println("counter:", counter)
}
