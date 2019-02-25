package redis_clients

import (
	"fmt"
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/RivenZoo/backbone/redis_connector"
	"github.com/RivenZoo/injectgo"
	"gopkg.in/redis.v5"
	"reflect"
)

// All redis client injected with its inject name like "redis_clients.test.0"
// db name "default" is also injected without inject name
// type Foo struct {
// 	 Client  	*redis.Client `inject:""` // get default redis db 0, same with "redis_clients.default.0"
//   Client1 	*redis.Client `inject:"redis_clients.1"` // get default redis db 1, same with "redis_clients.default.1"
//   TestClient *redis.Client `inject:"redis_clients.test.1"` // get test redis db 1
// }
//
// ** redis client json config format, more options see redis_connector/redis_connector.go#RedisOptions **
//{
//	"redis_config": {
//		"test": {
//			"addr": "127.0.0.1:6379",
//			"db_indexs": [0, 1, 2]
//		}
//	}
//}

var redisConnectorProviderObj *redisConnectorProvider

// use redisConnectorProvider to prevent *redis_connector.RedisConnector Close called by inject container.
type redisConnectorProvider struct {
	redisConnector *redis_connector.RedisConnector
}

func GetBootConfigKey() string {
	return "res.redis_clients"
}

func init() {
	key := GetBootConfigKey()
	c := objects_container.GetObjectContainer()
	c.ProvideFunc(injectgo.InjectFunc{
		Fn: func() (*redisConnectorProvider, error) {
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
			return &redisConnectorProvider{rc}, nil

		},
		Receiver: &redisConnectorProviderObj,
	})
}

func GetClient(name string, db int) *redis.Client {
	return redisConnectorProviderObj.redisConnector.GetClient(redis_connector.RedisDBOption{
		Name:    name,
		DBIndex: db,
	})
}

func GetInjectInfo(name string, db int) (injectName string, tp reflect.Type) {
	injectName = fmt.Sprintf("redis_clients.%s.%d", name, db)
	tp = reflect.TypeOf((*redis.Client)(nil))
	return
}

func injectRedisDBClients(rc *redis_connector.RedisConnector) {
	c := objects_container.GetObjectContainer()
	redisDBClients := rc.AllRedisDBClients()
	for i := range redisDBClients {
		cli := &redisDBClients[i]
		provideInjectRedisDBClient(c, cli.Name, cli.DBIndex, cli.Client)
	}
}

const defaultRedisName = "default"
const defaultRedisDBIndex = 0

func provideInjectRedisDBClient(c *injectgo.Container, name string, db int, cli *redis.Client) {
	injectName, _ := GetInjectInfo(name, db)
	c.ProvideByName(injectName, cli)

	if name == defaultRedisName {
		injectName = redisDBInjectNameWithoutName(db)
		c.ProvideByName(injectName, cli) // inject with name "redis_clients.{db}"

		if db == defaultRedisDBIndex {
			c.Provide(cli) // inject without name
		}
	}
}

func redisDBInjectNameWithoutName(db int) string {
	return fmt.Sprintf("redis_clients.%d", db)
}
