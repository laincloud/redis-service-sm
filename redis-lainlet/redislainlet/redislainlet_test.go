package redislainlet

import (
	"fmt"
	"testing"
	"time"
)

func Test_LainLet(t *testing.T) {
	go StartLainLet()
	for {
		time.Sleep(2 * time.Second)
		fmt.Println(Redis_addrs())
		fmt.Println(Sentinel_addrs())
	}
}
