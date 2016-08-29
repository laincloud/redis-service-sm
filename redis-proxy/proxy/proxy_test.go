package proxy

import (
	"math/rand"
	"strconv"
	"testing"

	redis "gopkg.in/redis.v4"
)

const (
	RedisAddr = "127.0.0.1:6379"
)

func thread_test(client *redis.Client) {
	for {
		client.Set(strconv.Itoa(rand.Intn(10000)), rand.Intn(10000), 0)
	}
}

func Test_redis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
	singal := make(chan struct{}, 1)
	for i := 0; i < 100; i++ {
		go thread_test(client)
	}
	<-singal
}
