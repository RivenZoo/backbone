package db_clients

import (
	"fmt"
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/resource_manager"
	"github.com/RivenZoo/backbone/resources"
	"github.com/RivenZoo/dsncfg"
	"github.com/RivenZoo/sqlagent"
	"github.com/jmoiron/sqlx/reflectx"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

type sqlConfig struct {
	ColumnMapTag string `json:"column_map_tag"`
}

type dbConfig struct {
	DB         dsncfg.Database         `json:"db"`
	Connection dsncfg.ConnectionConfig `json:"connection"`
	Sql        sqlConfig               `json:"sql"`
}

// DBConfig config json format
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

type namedSqlAgent struct {
	name     string
	sqlagent *sqlagent.SqlAgent
}

type namedSqlAgents struct {
	sqlagents []namedSqlAgent
	config    []dbConfig
}

func newNamedSqlAgent(cfg *DBConfig) (*namedSqlAgents, error) {
	ret := &namedSqlAgents{
		sqlagents: make([]namedSqlAgent, 0, len(cfg.NamedDBs)),
		config:    make([]dbConfig, 0, len(cfg.NamedDBs)),
	}
	for name, dbCfg := range cfg.NamedDBs {
		sa, err := sqlagent.NewSqlAgent(&dbCfg.DB)
		if err != nil {
			return nil, err
		}
		ret.sqlagents = append(ret.sqlagents, namedSqlAgent{
			name:     name,
			sqlagent: sa,
		})
		ret.config = append(ret.config, dbCfg)
	}
	return ret, nil
}

func setSqlConfig(sa *sqlagent.SqlAgent, cfg sqlConfig) {
	if cfg.ColumnMapTag != "" {
		sa.SetDBMapper(reflectx.NewMapperFunc(cfg.ColumnMapTag, strings.ToLower))
	}
}

func (na *namedSqlAgents) Init() error {
	for i := range na.config {
		na.sqlagents[i].sqlagent.SetConnectionConfig(na.config[i].Connection)
		setSqlConfig(na.sqlagents[i].sqlagent, na.config[i].Sql)
	}
	return nil
}

func (na *namedSqlAgents) Close() error {
	for i := range sqlagentContainer.sqlagents {
		if sqlagentContainer.sqlagents[i].sqlagent != nil {
			sqlagentContainer.sqlagents[i].sqlagent.Close()
		}
	}
	return nil
}

var sqlagentContainer *namedSqlAgents

func GetBootConfigKey() string {
	return "res.db_clients"
}

func GetClient(name string) *sqlagent.SqlAgent {
	for i := range sqlagentContainer.sqlagents {
		if sqlagentContainer.sqlagents[i].name == name {
			return sqlagentContainer.sqlagents[i].sqlagent
		}
	}
	return nil
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
		fmt.Println()
		return newNamedSqlAgent(cfg)
	}, &sqlagentContainer)
	resources.GetResourceContainer().RegisterCreator(key, creator)
}
