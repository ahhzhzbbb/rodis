package server

import (
	"fmt"
	"net"
)

const defaultPort = ":6379"

type Config struct {
	Port string
}

type Server struct {
	Config
	ln net.Listener
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
