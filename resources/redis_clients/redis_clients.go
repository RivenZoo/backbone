package redis_clients

import (
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/redis_connector"
	"github.com/RivenZoo/backbone/resource_manager"
	"github.com/RivenZoo/backbone/resources"
	"gopkg.in/redis.v5"
)

//
// ** Redis client json config format, more options see redis_connector/redis_connector.go#RedisOptions **
//{
//	"redis_config": {
//		"test": {
//			"addr": "127.0.0.1:6379",
//			"db_indexs": [0, 1, 2]
//		}
//	}
//}

var redisConnector *redis_connector.RedisConnector

func GetBootConfigKey() string {
	return "res.redis_clients"
}

func init() {
	key := GetBootConfigKey()
	creator := resource_manager.NewResourceCreator(func() (*redis_connector.RedisConnector, error) {
		g := bootconfig.GetConfigGetter()
		data, tp := g.GetConfig(key)
		var cfg *redis_connector.RedisConnectorConfig
		if err := data.Unmarshal(&cfg, tp); err != nil {
			return nil, err
		}
		rc, err := redis_connector.NewRedisConnector(*cfg)
		if err != nil {
			return nil, err
		}
		return rc, nil
	}, &redisConnector)
	resources.GetResourceContainer().RegisterCreator(key, creator)
}

func GetClient(name string, db int) *redis.Client {
	return redisConnector.GetClient(redis_connector.RedisDBOption{
		Name:    name,
		DBIndex: db,
	})
}
