package redis_clients

import (
	"fmt"
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/RivenZoo/backbone/redis_connector"
	"github.com/RivenZoo/backbone/resource_manager"
	"github.com/RivenZoo/backbone/resources"
	"gopkg.in/redis.v5"
	"reflect"
)

// redis client json config format
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
		injectRedisDBClients(rc)
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

func GetInjectInfo(name string, db int) (injectName string, tp reflect.Type) {
	injectName = fmt.Sprintf("res.redis_clients.%s.%d", name, db)
	tp = reflect.TypeOf((*redis.Client)(nil))
	return
}

func injectRedisDBClients(rc *redis_connector.RedisConnector) {
	c := objects_container.GetObjectContainer()
	redisDBClients := rc.AllRedisDBClients()
	for i := range redisDBClients {
		cli := &redisDBClients[i]
		injectName, _ := GetInjectInfo(cli.Name, cli.DBIndex)
		c.ProvideByName(injectName, cli.Client)
	}
}
