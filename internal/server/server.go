package server

import (
	"fmt"
	"net"
	"rodis/internal/engine"
	"time"
)

type Server struct {
	Config
	ln net.Listener
	kv *engine.KeyValue
}

func NewServer(cfg Config) *Server {
	if len(cfg.Port) == 0 {
		cfg.Port = defaultPort
	}
	return &Server{
		Config: cfg,
		kv:     engine.NewKeyValue(),
	}
}

func (s *Server) Start() bool {
	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		return false
	}
	s.ln = ln

	go s.runActiveExpiration()

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
		// fmt.Println("connection estalished!")
		go s.handleConnection(conn)
	}
}

func (s *Server) runActiveExpiration() {
	for {
		time.Sleep(time.Duration(s.Expire.CycleIntervalMs) * time.Millisecond)
		s.kv.ActiveExpiration(s.Expire.SampleSize, s.Expire.ExpireThreshold, s.Expire.TimeBudgetMs)
	}
}
