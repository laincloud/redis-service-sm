package proxy

import (
	"github.com/robfig/config"
	"os"
)

var (
	MASTER_NAME_SENTINEL = "master_service"
	CHECK_INTERVAL_MS    = 10000

	DEBUG = false
)

var (
	Max_client int
	Time_out   int
	Port       int

	ConnTimeoutSec int
	BufferSize     int
)

func Config_proxy(c *config.Config) {
	Max_client = 1000
	Time_out = 10000
	Port = 6379

	ConnTimeoutSec = 5
	BufferSize = 1024

	if c == nil {
		return
	}

	var err error

	if debug, err := c.Bool("proxy", "debug"); err == nil {
		DEBUG = debug
	}

	if check_interval_ms, err := c.Int("proxy", "check_interval_ms"); err == nil {
		CHECK_INTERVAL_MS = check_interval_ms
	}

	if master_name_sentinel, err := c.String("proxy", "master_name_sentinel"); err == nil {
		MASTER_NAME_SENTINEL = master_name_sentinel
	}

	if Max_client, err = c.Int("proxy", "max-client"); err != nil {
		Max_client = 1000
	}

	if Time_out, err = c.Int("proxy", "timeout"); err != nil {
		Time_out = 10000
	}

	if Port, err = c.Int("proxy", "port"); err != nil {
		Port = 8889
	}

	if BufferSize, err = c.Int("proxy", "buffersize"); err != nil {
		BufferSize = 1024
	}

	if ConnTimeoutSec, err = c.Int("proxy", "redisConTimeoutSec"); err != nil {
		ConnTimeoutSec = 5
	}

}

func Load_config(file_name string) {
	c, _ := config.ReadDefault(file_name)
	Config_proxy(c)
}

func GetEnvDefault(key, defaultValue string) string {
	s := os.Getenv(key)
	if s != "" {
		return s
	}
	return defaultValue
}
