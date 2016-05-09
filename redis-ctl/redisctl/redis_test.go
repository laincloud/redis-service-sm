package redisctl

import (
	"testing"
)

func Test_ListNodes(t *testing.T) {

	StartRedisMonitor()
	// sentinel_addrs := []*Address{BuildAddress("127.0.0.1", "5001"), BuildAddress("127.0.0.1", "5002"), BuildAddress("127.0.0.1", "5004")}
	// if err := RemoveSlave("127.0.0.1", "6003", "mymaster", sentinel_addrs...); err == nil {
	// 	t.Log("remove slave from sentinel pass")
	// }
}
