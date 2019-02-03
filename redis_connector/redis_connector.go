package redis_connector

import (
	"errors"
	"gopkg.in/redis.v5"
	"time"
)

const (
	maxDBIndex = 15
)

var (
	errOverMaxDBIndex = errors.New("over max db index")
	errUnknownRedis   = errors.New("no such name redis")
)

type RedisConnector struct {
	namedClients []*namedRedisClients
}

type namedRedisClients struct {
	name    string
	clients []*redis.Client
}

type RedisOptions struct {
	// host:port address.
	Addr string `json:"addr"`

	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string `json:"password"`

	DBIndexs []int `json:"db_indexs"`

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeoutMS int `json:"dial_timeout_ms"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is 3 seconds.
	ReadTimeoutMS int `json:"read_timeout_ms"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is 3 seconds.
	WriteTimeoutMS int `json:"write_timeout_ms"`

	// Maximum number of socket connections.
	// Default is 10 connections.
	PoolSize int `json:"pool_size"`
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeoutMS int `json:"pool_timeout_ms"`
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is to not close idle connections.
	IdleTimeoutMS int `json:"idle_timeout_ms"`

	// Enables read only queries on slave nodes.
	ReadOnly bool `json:"read_only"`
}

type RedisConnectorConfig struct {
	RedisConfig map[string]RedisOptions `json:"redis_config"`
}

func NewRedisConnector(connectorConfig RedisConnectorConfig) (*RedisConnector, error) {
	if len(connectorConfig.RedisConfig) == 0 {
		return nil, errors.New("no redis config")
	}
	ret := &RedisConnector{
		namedClients: make([]*namedRedisClients, 0, len(connectorConfig.RedisConfig)),
	}
	for name, opt := range connectorConfig.RedisConfig {
		if name == "" {
			return nil, errors.New("no redis name")
		}
		max := -1
		for _, idx := range opt.DBIndexs {
			if idx > max {
				max = idx
			}
		}
		if max > maxDBIndex {
			return nil, errOverMaxDBIndex
		}

		redisOpt := &redis.Options{
			Addr:         opt.Addr,
			Password:     opt.Password,
			DialTimeout:  time.Duration(opt.DialTimeoutMS) * time.Millisecond,
			ReadTimeout:  time.Duration(opt.ReadTimeoutMS) * time.Millisecond,
			WriteTimeout: time.Duration(opt.WriteTimeoutMS) * time.Millisecond,
			PoolSize:     opt.PoolSize,
			PoolTimeout:  time.Duration(opt.PoolTimeoutMS) * time.Millisecond,
			IdleTimeout:  time.Duration(opt.IdleTimeoutMS) * time.Millisecond,
			ReadOnly:     opt.ReadOnly,
		}
		clients := make([]*redis.Client, max+1)
		for _, idx := range opt.DBIndexs {
			redisOpt.DB = idx
			clients[idx] = redis.NewClient(redisOpt)
		}
		ret.namedClients = append(ret.namedClients, &namedRedisClients{
			name:    name,
			clients: clients,
		})
	}
	return ret, nil
}

type RedisDBOption struct {
	Name    string `json:"name"`
	DBIndex int    `json:"db_index"`
}

// GetClient get redis by name and db.
// If no redis or db will return nil.
func (c *RedisConnector) GetClient(opt RedisDBOption) *redis.Client {
	for i := range c.namedClients {
		if c.namedClients[i].name == opt.Name {
			return c.namedClients[i].clients[opt.DBIndex]
		}
	}
	return nil
}

func (c *RedisConnector) Close() error {
	for i := range c.namedClients {
		for _, cli := range c.namedClients[i].clients {
			if cli != nil {
				cli.Close()
			}
		}
	}
	return nil
}
