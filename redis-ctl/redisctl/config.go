package redisctl

import (
	"github.com/mijia/sweb/log"
	"github.com/robfig/config"
)

var (
	REDIS_MANAGER_INTERVAL_MS = 10000
	MASTER_NAME_SENTINEL      = "master_service"

	DEBUG = false
)

func ConfigCtl(file_name string) error {
	c, err := config.ReadDefault(file_name)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if debug, err := c.Bool("ctl", "debug"); err == nil {
		DEBUG = debug
	}
	if redis_check_interval_ms, err := c.Int("ctl", "redis_check_interval_ms"); err == nil {
		REDIS_MANAGER_INTERVAL_MS = redis_check_interval_ms
	}

	if master_name_sentinel, err := c.String("monitor", "master_name_sentinel"); err == nil {
		MASTER_NAME_SENTINEL = master_name_sentinel
	}
	return nil
}
