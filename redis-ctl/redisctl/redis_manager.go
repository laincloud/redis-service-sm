package redisctl

import (
	"github.com/laincloud/redis-libs/redislibs"
	"github.com/mijia/sweb/log"
	"strconv"
)

func CheckAndInitRedisSentinel(sentinel_addrs map[int]*redislibs.Address, master *redislibs.Address) {
	// SENTINEL MONITOR <name> <ip> <port> <quorum>
	if master == nil {
		return
	}
	quorum := len(sentinel_addrs)/2 + 1
	for _, s_addr := range sentinel_addrs {
		monitorAndCfgSentinel(s_addr, master, quorum)
	}
}

func monitorAndCfgSentinel(s_addr *redislibs.Address, master *redislibs.Address, quorum int) {
	if master_addr, err := redislibs.GetMasterAddrByName(s_addr.Host, s_addr.Port, MASTER_NAME_SENTINEL); err == nil && master_addr == nil {
		redislibs.MonitorSentinel(s_addr.Host, s_addr.Port, master.Host, master.Port, MASTER_NAME_SENTINEL, strconv.Itoa(quorum))
		redislibs.ConfigSentinel(s_addr.Host, s_addr.Port, MASTER_NAME_SENTINEL, "down-after-milliseconds", "30000")
		redislibs.ConfigSentinel(s_addr.Host, s_addr.Port, MASTER_NAME_SENTINEL, "parallel-syncs", "1")
		redislibs.ConfigSentinel(s_addr.Host, s_addr.Port, MASTER_NAME_SENTINEL, "failover-timeout", "180000")
	}
}

func checkAndFixRoleStatus(sys_master *redislibs.Address) {
	if role, status, err := redislibs.RoleStatus(sys_master.Host, sys_master.Port); err == nil {
		if role == redislibs.ROLE_SLAVE && status == redislibs.ROLE_STATUS_CONNECT {
			log.Warn("shit problem happend")
			redislibs.SlaveOf(sys_master.Host, sys_master.Port, "no", "one")
		}
	}
}

func getSlaves(sentinel_addrs map[int]*redislibs.Address) []*redislibs.Address {
	for _, s_addr := range sentinel_addrs {
		if slaves, err := redislibs.GetSlavesInSentinel(s_addr.Host, s_addr.Port, MASTER_NAME_SENTINEL); err == nil {
			return slaves
		}
	}
	return make([]*redislibs.Address, 0, 0)
}

func getMaster(sentinel_addrs map[int]*redislibs.Address) *redislibs.Address {
	for _, s_addr := range sentinel_addrs {
		if master_addr, err := redislibs.GetMasterAddrByName(s_addr.Host, s_addr.Port, MASTER_NAME_SENTINEL); err == nil {
			return master_addr
		}
	}
	return nil
}

func ManagerRedisService(redis_addrs map[int]*redislibs.Address, sentinel_addrs map[int]*redislibs.Address) {
	if len(redis_addrs) <= 1 || len(sentinel_addrs) == 0 {
		return
	}
	if sys_master == nil {
		sys_master = redis_addrs[1]
	}
	CheckAndInitRedisSentinel(sentinel_addrs, sys_master)
	sys_master = getMaster(sentinel_addrs)
	if sys_master == nil {
		log.Warn("No Master")
		return
	}
	slaves := getSlaves(sentinel_addrs)
	log.Debug("sys_master: ", sys_master)
	in_cluster := false
	for _, r_addr := range redis_addrs {
		in_cluster = false
		if r_addr.Host == sys_master.Host && r_addr.Port == sys_master.Port {
			continue
		}
		for _, s_addr := range slaves {
			if s_addr.Host == r_addr.Host && s_addr.Port == r_addr.Port {
				in_cluster = true
				break
			}
		}
		if !in_cluster {
			redislibs.SlaveOf(r_addr.Host, r_addr.Port, sys_master.Host, sys_master.Port)
		}
	}
	checkAndFixRoleStatus(sys_master)
}
