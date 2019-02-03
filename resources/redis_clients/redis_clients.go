package redis_clients

import (
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/redis_connector"
	"github.com/RivenZoo/backbone/resource_manager"
	"github.com/RivenZoo/backbone/resources"
	"gopkg.in/redis.v5"
)

var redisConnector *redis_connector.RedisConnector

func GetBootConfigKey() string {
	return "res.redis_clients"
}

func init() {
	key := GetBootConfigKey()
	creator := resource_manager.NewResourceCreator(func() (*redis_connector.RedisConnector, error) {
		g := bootconfig.GetConfigGetter()
		data := g.GetConfig(key)
		var cfg *redis_connector.RedisConnectorConfig
		if err := data.Unmarshal(&cfg); err != nil {
			return nil, err
		}
		return redis_connector.NewRedisConnector(*cfg)
	}, &redisConnector)
	resources.GetResourceContainer().RegisterCreator(key, creator)
}

func GetClient(name string, db int) *redis.Client {
	return redisConnector.GetClient(redis_connector.RedisDBOption{
		Name:    name,
		DBIndex: db,
	})
}
