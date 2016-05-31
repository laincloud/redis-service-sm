package proxy

import (
	// "fmt"
	// "runtime"
	// "os"
	"github.com/mijia/sweb/log"
	"testing"
	"time"
)

// func Test_LainLet(t *testing.T) {
// 	// LainLet()
// 	fmt.Println(runtime.GOOS)
// }

func Test_Kqueue(t *testing.T) {
	defer log.Info("out")
	Load_config("../proxy.conf")
	if DEBUG {
		log.EnableDebug()
	}
	for {
		p := NewProxy()
		p.StartServer()
		p.StopServer()
		time.Sleep(60 * time.Second)
	}

}
