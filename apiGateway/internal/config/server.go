package config

import "time"

type Server struct {
	ListenHost         string `yaml:"listenHost"`
	ListenPort         string `yaml:"listenPort"`
	Scheme             string `yaml:"scheme"`
	ShutdownTimeoutSec int    `yaml:"shutdownTimeoutSec"`
}

func (s *Server) ShutdownTimeout() time.Duration {
	return time.Duration(s.ShutdownTimeoutSec) * time.Second
}

func (s *Server) Addr() string {
	if s.ListenPort == "" {
		return s.ListenHost
	}

	return s.ListenHost + ":" + s.ListenPort
}
