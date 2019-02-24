package model

import (
	"github.com/RivenZoo/backbone/objects_container"
	"gopkg.in/redis.v5"
)

var counter = &redisCounter{}

type redisCounter struct {
	Client *redis.Client `inject:"res.redis_clients.test_redis.1"`
}

func (c *redisCounter) Count(key string) (int64, error) {
	return c.Client.Incr(key).Result()
}

func (c *redisCounter) Get(key string) (int64, error) {
	return c.Client.Get(key).Int64()
}

func GetRedisCounter() *redisCounter {
	return counter
}

func init() {
	objC := objects_container.GetObjectContainer()
	objC.Provide(counter)
}
