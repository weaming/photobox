package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

var codec *cache.Codec

func initCodec() {
	_, codec = NewRedisCache(map[string]string{"server0": redisHostPort}, redisDB)
}

func NewRedisCache(hostPort map[string]string, db int) (*redis.Ring, *cache.Codec) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: hostPort,
		DB:    db,
	})

	codec := &cache.Codec{
		Redis: ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return json.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return json.Unmarshal(b, v)
		},
	}
	return ring, codec
}

func CacheSet(key string, obj interface{}) error {
	return codec.Set(&cache.Item{
		Key:        key,
		Object:     obj,
		Expiration: time.Hour * 24 * 365,
	})
}

func CacheGet(key string, obj interface{}) error {
	return codec.Get(key, obj)
}
