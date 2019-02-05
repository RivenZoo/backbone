package model

import (
	"github.com/RivenZoo/backbone/resources/redis_clients"
	"time"
)

func SetAbbreviateURL(abbrURL, srcURL string) error {
	cli := redis_clients.GetClient("test_redis", 0)
	return cli.Set(abbrURL, srcURL, time.Hour*7*24).Err()
}

func GetSourceURL(abbrURL string) (string, error) {
	cli := redis_clients.GetClient("test_redis", 0)
	return cli.Get(abbrURL).Result()
}
