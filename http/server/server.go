package server

import (
	"context"
	"github.com/RivenZoo/backbone/http/logger"
	"net/http"
	"time"
)

const defaultShutdownTimeoutMS = 10000

type ServerConfig struct {
	Addr                string `json:"addr"`
	ReadTimeoutMS       int    `json:"read_timeout_ms"`
	ReadHeaderTimeoutMS int    `json:"read_header_timeout_ms"`
	WriteTimeoutMS      int    `json:"write_timeout_ms"`
	IdleTimeoutMS       int    `json:"idle_timeout_ms"`
	MaxHeaderBytes      int    `json:"max_header_bytes"`
	ShutdownTimeoutMS   int    `json:"shutdown_timeout_ms"`
}

type SimpleServer struct {
	server            *http.Server
	shutdownTimeoutMS int
}

func NewSimpleServer(cfg *ServerConfig) (*SimpleServer, error) {
	ret := &SimpleServer{
		server: &http.Server{
			Addr:              cfg.Addr,
			ReadTimeout:       time.Duration(cfg.ReadTimeoutMS) * time.Millisecond,
			ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeoutMS) * time.Millisecond,
			WriteTimeout:      time.Duration(cfg.WriteTimeoutMS) * time.Millisecond,
			IdleTimeout:       time.Duration(cfg.IdleTimeoutMS) * time.Millisecond,
			MaxHeaderBytes:    cfg.MaxHeaderBytes,
		},
		shutdownTimeoutMS: cfg.ShutdownTimeoutMS,
	}
	if ret.shutdownTimeoutMS <= 0 {
		ret.shutdownTimeoutMS = defaultShutdownTimeoutMS
	}
	return ret, nil
}

func (s *SimpleServer) SetHTTPHandler(h http.Handler) {
	s.server.Handler = h
}

func (s *SimpleServer) Run() error {
	if s.server.Handler == nil {
		s.server.Handler = helloHandler{"welcome to simple server!"}
	}
	s.server.SetKeepAlivesEnabled(true)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logf("[ERROR] ListenAndServe error %v", err)
		return err
	}
	return nil
}

func (s *SimpleServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.shutdownTimeoutMS)*time.Millisecond)
	defer cancel()

	return s.server.Shutdown(ctx)
}
