package db_clients

import (
	"fmt"
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/resource_manager"
	"github.com/RivenZoo/backbone/resources"
	"github.com/RivenZoo/dsncfg"
	"github.com/RivenZoo/sqlagent"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

var sqlagentContainer *namedSqlAgents

func GetBootConfigKey() string {
	return "res.db_clients"
}

type sqlConfig struct {
	ColumnMapTag string `json:"column_map_tag"`
}

type dbConfig struct {
	DB         dsncfg.Database         `json:"db"`
	Connection dsncfg.ConnectionConfig `json:"connection"`
	Sql        sqlConfig               `json:"sql"`
}

// ** DBConfig config json format **
// {
// "dbs": {
//		"test-db": {
//			"db": {
//				"host": "127.0.0.1",
//				"port": 3306,
//				"name": "test",
//				"type": "mysql",
//				"user": "test",
//				"password": "",
//				"parameters": {
//				"parseTime": "true",
//				"charset": "utf8mb4,utf8",
//				"loc": "Asia/Shanghai"
//				},
//			},
//			"connection": {
//				"max_open_connections": 60,
//				"max_idle_connections": 40,
//				"max_life_time": 600
//			},
//			"sql": {
//				"column_map_tag": "json"
//			}
//		}
//	}
// }
type DBConfig struct {
	NamedDBs map[string]dbConfig `json:"dbs"`
}

func GetClient(name string) *sqlagent.SqlAgent {
	for i := range sqlagentContainer.sqlagents {
		if sqlagentContainer.sqlagents[i].name == name {
			return sqlagentContainer.sqlagents[i].sqlagent
		}
	}
	return nil
}

func GetInjectInfo(name string) (injectName string, tp reflect.Type) {
	injectName = fmt.Sprintf("res.db_clients.%s", name)
	tp = reflect.TypeOf((*sqlagent.SqlAgent)(nil))
	return
}

func init() {
	key := GetBootConfigKey()
	creator := resource_manager.NewResourceCreator(func() (*namedSqlAgents, error) {
		g := bootconfig.GetConfigGetter()
		data, tp := g.GetConfig(key)
		var cfg *DBConfig
		if err := data.Unmarshal(&cfg, tp); err != nil {
			return nil, err
		}

		namedSA, err := newNamedSqlAgent(cfg)
		if err != nil {
			return nil, err
		}
		return namedSA, nil
	}, &sqlagentContainer)
	resources.GetResourceContainer().RegisterCreator(key, creator)
}
