package redislainlet

import (
	"os"
)

var (
	LAINLET_ADDR = GetEnvDefault("LAINLET_ADDR", "lainlet.lain:9001")
	APPNAME      = GetEnvDefault("LAIN_APPNAME", "redis-service-sm")
	WATHCER      = "/v2/procwatcher/?appname=" + APPNAME + "&heartbeat=5"

	PROC_REDIS_NAME   = GetEnvDefault("PROC_REDIS_NAME", "redis")
	PROC_SENTILE_NAME = GetEnvDefault("PROC_SENTILE_NAME", "redis-sentinel")

	ERROR_IDLE_TIME = 10000
)

func GetEnvDefault(key, defaultValue string) string {
	s := os.Getenv(key)
	if s != "" {
		return s
	}
	return defaultValue
}
