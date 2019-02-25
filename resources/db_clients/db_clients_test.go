package db_clients

import (
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/configutils"
	"github.com/RivenZoo/backbone/resources"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type dbBootConfig struct {
}

func (cfg dbBootConfig) GetConfig(key string) (bootconfig.RawConfigData, configutils.ConfigType) {
	if key == GetBootConfigKey() {
		return []byte(`
{
"dbs": {
		"test-db": {
			"db": {
				"host": "127.0.0.1",
				"port": 3306,
				"name": "test",
				"type": "mysql",
				"user": "root",
				"password": "testadmin",
				"parameters": {
				"parseTime": "true",
				"charset": "utf8mb4,utf8",
				"loc": "Asia/Shanghai"
				}
			},
			"connection": {
				"max_open_connections": 60,
				"max_idle_connections": 40,
				"max_life_time": 600
			},
			"sql": {
				"column_map_tag": "json"
			}
		}
	}
}
`), configutils.ConfigTypeJSON
	}
	return []byte{}, ""
}

func TestMain(m *testing.M) {
	// init config first
	bootconfig.RegisterConfigGetter(dbBootConfig{})

	rc := resources.GetResourceContainer()
	rc.Init()
	defer rc.Close()

	os.Exit(m.Run())
}

func TestGetClient(t *testing.T) {
	rc := resources.GetResourceContainer()
	obj := rc.GetResource(GetBootConfigKey())
	assert.NotNil(t, obj)

	cli := GetClient("test-db")
	assert.NotNil(t, cli)
}
