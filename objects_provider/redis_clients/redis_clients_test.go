package redis_clients

import (
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/configutils"
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/stretchr/testify/assert"
	"gopkg.in/redis.v5"
	"os"
	"testing"
	"unsafe"
)

type redisBootConfig struct {
}

func (cfg redisBootConfig) GetConfig(key string) (bootconfig.RawConfigData, configutils.ConfigType) {
	if key == GetBootConfigKey() {
		return []byte(`
{
  "redis_config": {
	 "default": {
	   "addr": "127.0.0.1:6379",
 	   "db_indexs": [0, 1, 2]
	 },
	"test": {
	   "addr": "127.0.0.1:6379",
 	   "db_indexs": [0, 1, 2]
	 }
  }
}
`), configutils.ConfigTypeJSON
	}
	return []byte{}, ""
}

func TestMain(m *testing.M) {
	// set config first
	bootconfig.RegisterConfigGetter(redisBootConfig{})

	os.Exit(m.Run())
}

type Foo struct {
	DefaultCli *redis.Client `inject:""`
	Cli1       *redis.Client `inject:"redis_clients.1"`
	TestCli    *redis.Client `inject:"redis_clients.test.0"`
}

func TestRedisClientInject(t *testing.T) {
	foo := &Foo{}
	c := objects_container.GetObjectContainer()
	c.Provide(foo)
	objects_container.Init()
	defer objects_container.Close()

	assert.NotNil(t, foo.DefaultCli)
	assert.NotNil(t, foo.Cli1)
	assert.NotNil(t, foo.TestCli)

	t.Log(foo.DefaultCli.Ping().Result())

	cli := GetClient("default", 0)
	assert.Equal(t, unsafe.Pointer(cli), unsafe.Pointer(foo.DefaultCli))

	cli = GetClient("default", 1)
	assert.Equal(t, unsafe.Pointer(cli), unsafe.Pointer(foo.Cli1))

	cli = GetClient("test", 0)
	assert.Equal(t, unsafe.Pointer(cli), unsafe.Pointer(foo.TestCli))
}
