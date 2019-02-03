package redis_clients

import (
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/resources"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type redisBootConfig struct {
}

func (cfg redisBootConfig) GetConfig(key string) bootconfig.RawConfigData {
	if key == GetBootConfigKey() {
		return []byte(`
{
  "redis_config": {
	 "test": {
	   "addr": "127.0.0.1:6379",
 	   "db_indexs": [0, 1, 2]
	 }
  }
}
`)
	}
	return []byte{}
}

func TestMain(m *testing.M) {
	// set config first
	bootconfig.RegisterConfigGetter(redisBootConfig{})
	// init resource manager
	c := resources.GetResourceContainer()
	c.Init()
	defer c.Close()
	os.Exit(m.Run())
}

func TestGetClient(t *testing.T) {
	cli := GetClient("test", 0)
	assert.NotNil(t, cli)
}
