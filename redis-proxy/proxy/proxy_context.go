package proxy

import (
	"github.com/laincloud/redis-libs/redislibs"
)

var (
	master_addr         *redislibs.Address
	avail_sentinel_addr *redislibs.Address
)
