package db_clients

import (
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/configutils"
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/RivenZoo/sqlagent"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"unsafe"
)

type dbBootConfig struct {
}

func (cfg dbBootConfig) GetConfig(key string) (bootconfig.RawConfigData, configutils.ConfigType) {
	if key == GetBootConfigKey() {
		return []byte(`
{
"dbs": {
		"default": {
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
		},
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

	os.Exit(m.Run())
}

type Foo struct {
	DBClient  *sqlagent.SqlAgent `inject:""`
	TestDBCli *sqlagent.SqlAgent `inject:"db_clients.test-db"`
}

func TestInjectDBClient(t *testing.T) {
	foo := &Foo{}
	objects_container.GetObjectContainer().Provide(foo)

	objects_container.Init()
	defer objects_container.Close()

	cli := GetClient("test-db")
	assert.NotNil(t, cli)

	assert.NotNil(t, foo.DBClient)
	assert.NotNil(t, foo.TestDBCli)
	assert.Equal(t, unsafe.Pointer(cli), unsafe.Pointer(foo.TestDBCli))
}
