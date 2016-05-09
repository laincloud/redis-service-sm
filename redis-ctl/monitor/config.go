package monitor

import (
	"github.com/mijia/sweb/log"
	"github.com/robfig/config"
	"os"
	"strings"
)

var (
	quotas     = make(map[string][]string)
	DOMAIN     = strings.Replace(GetEnvDefault("LAIN_DOMAIN", "lain.local"), ".", "_", -1)
	APPNAME    = strings.Replace(GetEnvDefault("LAIN_APPNAME", "redis_service_sm"), "-", "_", -1)
	KEY_PREFIX = DOMAIN + "." + strings.TrimPrefix(APPNAME, "resource.")

	DEBUG = false

	GRAPHITE_DOMAIN = "graphite.lain"
	GRAPHITE_PORT   = "2003"

	REDIS_MONITOR_INTERVAL_MS = 60000
)

func ConfigMonitor(file_name string) error {
	c, err := config.ReadDefault(file_name)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if host, err := c.String("graphite", "host"); err == nil {
		GRAPHITE_DOMAIN = host
	}
	if port, err := c.String("graphite", "port"); err == nil {
		GRAPHITE_PORT = port
	}

	if debug, err := c.Bool("monitor", "debug"); err == nil {
		DEBUG = debug
	}
	if interval_ms, err := c.Int("monitor", "monitor_interval_ms"); err == nil {
		REDIS_MONITOR_INTERVAL_MS = interval_ms
	}

	if m_qstr, err := c.String("monitor", "quota"); err == nil {
		m_qs := strings.Split(m_qstr, ",")
		for _, m_q := range m_qs {
			if s_qstr, err := c.String("metrics", m_q); err == nil {
				s_q := strings.Split(s_qstr, ",")
				quotas[m_q] = s_q
			} else {
				return err
			}
		}
	}
	return nil
}

func GetEnvDefault(key, defaultValue string) string {
	s := os.Getenv(key)
	if s != "" {
		return s
	}
	return defaultValue
}
