package server

import (
	"fmt"
	"net"
	"rodis/internal/engine"
)

type Server struct {
	Config
	ln net.Listener
	kv *engine.KeyValue
	et *engine.ExpireTime
}

func NewServer(cfg Config) *Server {
	if len(cfg.Port) == 0 {
		cfg.Port = defaultPort
	}
	return &Server{
		Config: cfg,
		kv:     engine.NewKeyValue(),
		et:     engine.NewExpireTime(),
	}
}

func (s *Server) Start() bool {
	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		return false
	}
	s.ln = ln
	s.loop()
	return true
}

func (s *Server) loop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("failed to accept connection, err:", err)
			continue
		}
		fmt.Println("connection estalished!")
		go s.handleConnection(conn)
	}
}
