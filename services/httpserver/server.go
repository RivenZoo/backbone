package httpserver

import (
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/http/server"
	"github.com/RivenZoo/backbone/service_manager"
	"github.com/RivenZoo/backbone/services"
)

var defaultServer *server.SimpleServer

// server config json format
// {
//    "addr": "127.0.0.1:8080",
//    "read_timeout_ms": 0,
//    "read_header_timeout_ms": 0,
//    "write_timeout_ms": 0,
//    "idle_timeout_ms": 0,
//    "max_header_bytes": 0,
//    "shutdown_timeout_ms": 0
// }

func GetHTTPServer() *server.SimpleServer {
	return defaultServer
}

func GetBootConfigKey() string {
	return "svc.httpserver"
}

func init() {
	key := GetBootConfigKey()
	creator := service_manager.NewServiceCreator(func() (*server.SimpleServer, error) {
		g := bootconfig.GetConfigGetter()
		data, tp := g.GetConfig(key)
		var cfg *server.ServerConfig
		if err := data.Unmarshal(&cfg, tp); err != nil {
			return nil, err
		}
		return server.NewSimpleServer(cfg)
	}, &defaultServer)
	services.GetServiceContainer().RegisterCreator(key, creator)
}
