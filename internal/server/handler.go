package server

import (
	"fmt"
	"io"
	"net"
	"rodis/internal/command"
	"rodis/internal/protocol/resp"
)

const defaultPort = ":6379"

type Config struct {
	Port string
}

type Server struct {
	Config
	ln net.Listener
}

func NewServer(cfg Config) *Server {
	if len(cfg.Port) == 0 {
		cfg.Port = defaultPort
	}
	return &Server{
		Config: cfg,
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	rp := resp.NewResp(conn)

	for {
		request, err := rp.ParseRESP()
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to parse clients request")
			} else {
				fmt.Printf("client %s disconnected\n", conn.RemoteAddr())
			}
			s.removeConnection(conn)
			return
		}
		fmt.Printf("request: %v\n", request)
		response := command.Commands(request)
		encoder := response.Marshal()
		rp.Writer(encoder)
	}
}

func (s *Server) removeConnection(conn net.Conn) {

}

// func (s *Server) handleCommand(command string) {
// 	switch command {
// 	case "PING":
// 		v := ping()
// 	}
// }
