package db_clients

import (
	"fmt"
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/objects_container"
	"github.com/RivenZoo/dsncfg"
	"github.com/RivenZoo/injectgo"
	"github.com/RivenZoo/sqlagent"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

var sqlagentProviderObj *sqlAgentProvider

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

// All db injected with its inject name like "db_clients.test-db" to objects_container.
// db name "default" is also injected without inject name
// type Foo struct {
// 	 SqlAgent 		*sqlagent.SqlAgent `inject:""` // get default db, same with `inject:"db_clients.default"`
// 	 TestSqlAgent 	*sqlagent.SqlAgent `inject:"db_clients.test-db"` // get test-db db
// }
//
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
	for i := range sqlagentProviderObj.namedSqlAgents.sqlagents {
		if sqlagentProviderObj.namedSqlAgents.sqlagents[i].name == name {
			return sqlagentProviderObj.namedSqlAgents.sqlagents[i].sqlagent
		}
	}
	return nil
}

func GetInjectInfo(name string) (injectName string, tp reflect.Type) {
	injectName = fmt.Sprintf("db_clients.%s", name)
	tp = reflect.TypeOf((*sqlagent.SqlAgent)(nil))
	return
}

// use sqlAgentProvider to prevent *namedSqlAgents Close called by inject container.
type sqlAgentProvider struct {
	namedSqlAgents *namedSqlAgents
}

func init() {
	key := GetBootConfigKey()
	c := objects_container.GetObjectContainer()
	c.ProvideFunc(injectgo.InjectFunc{
		Fn: func() (*sqlAgentProvider, error) {
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
			// call Init manually
			if err = namedSA.Init(); err != nil {
				return nil, err
			}

			injectNamedSqlAgents(namedSA)
			return &sqlAgentProvider{namedSA}, nil
		},
		Receiver: &sqlagentProviderObj,
	})
}

func injectNamedSqlAgents(namedSA *namedSqlAgents) {
	c := objects_container.GetObjectContainer()
	for i := range namedSA.sqlagents {
		sa := &namedSA.sqlagents[i]
		provideInjectSqlAgent(c, sa.name, sa.sqlagent)
	}
}

const defaultDBName = "default"

func provideInjectSqlAgent(c *injectgo.Container, name string, sa *sqlagent.SqlAgent) {
	injectName, _ := GetInjectInfo(name)
	c.ProvideByName(injectName, sa) // inject with name db_clients.{dbName}

	// set default inject, use injected SqlAgent without name
	// eg.
	// type A struct
	// {
	// 	SqlAgent *sqlagent.SqlAgent `inject:""`
	// }
	if name == defaultDBName {
		c.Provide(sa)
	}
}
