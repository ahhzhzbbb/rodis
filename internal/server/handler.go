package server

import (
	"fmt"
	"io"
	"net"
	"rodis/internal/command"
	"rodis/internal/factory"
	"rodis/internal/protocol/resp"
	"strings"
)

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
				return
			}
			// s.removeConnection(conn)
			continue
		}
		fmt.Printf("request: %v\n", request)

		if len(request.Array) == 0 {
			continue
		}

		response := s.handleRequest(request)
		fmt.Printf("response: %v\n", response)

		rp.WriteBytes(response)
	}
}

func (s *Server) handleRequest(request resp.Value) resp.Value {
	var result resp.Value
	typeOfCommand := strings.ToUpper(request.Array[0].Bulk)

	comCreator, ok := factory.CommandRegistry[typeOfCommand]
	if !ok {
		result = resp.NewError(fmt.Sprintf("ERR unknown command '%s', with args beginning with: ", typeOfCommand))
	} else {
		comm := comCreator()
		result = comm.Execute(request.Array[1:], command.NewCommandContext(s.kv, s.et))
	}
	return result
}

func (s *Server) removeConnection(conn net.Conn) {
	//removing connection from registry...
}
