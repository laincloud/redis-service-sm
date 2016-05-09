package monitor

import (
	"fmt"
	"github.com/laincloud/redis-libs/redislibs"
)

const (
	KEY_FORMAT = "%s.%d.%s.%s" // domain.appname.instanceNo.Quota.Param
)

func RedisServerMetrics(instanceNo int, host, port string) (map[string]string, error) {
	infos, err := redislibs.RedisNodeInfo(host, port)
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	var title string
	for k, qs := range quotas {
		title = k
		for _, q := range qs {
			key := fmt.Sprintf(KEY_FORMAT, KEY_PREFIX, instanceNo, title, q) //ENDPOINT+"."+strconv.Itoa(instanceNo)+"."+title + "." + q
			res[key] = infos[title][q]
		}
	}
	return res, nil
}
