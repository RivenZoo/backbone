package db_clients

import (
	"github.com/RivenZoo/sqlagent"
	"github.com/jmoiron/sqlx/reflectx"
	"strings"
)

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
	for i := range na.sqlagents {
		if na.sqlagents[i].sqlagent != nil {
			na.sqlagents[i].sqlagent.Close()
		}
	}
	return nil
}