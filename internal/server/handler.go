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
		var response resp.Value
		decoder, err := rp.ParseRESP()
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to parse clients request")
			} else {
				fmt.Printf("client %s disconnected\n", conn.RemoteAddr())
			}
			fmt.Println("Error: ", err)
			s.removeConnection(conn)
			return
		}
		fmt.Printf("request: %v\n", decoder)
		if len(decoder.Array) == 0 {
			continue
		}
		typeOfCommand := decoder.Array[0].Bulk

		typeOfCommand = strings.ToUpper(typeOfCommand)

		creator, ok := factory.CommandRegistry[typeOfCommand]
		if !ok {
			response = resp.NewError("FAILED")
		} else {
			comm := creator()
			response = comm.Execute(decoder.Array[1:], command.NewCommandContext(s.kv, s.et))
		}

		fmt.Printf("response: %v\n", response)

		encoder := response.Marshal()
		err = rp.Writer(encoder)
		if err != nil {
			//do something
		}
	}
}

func (s *Server) removeConnection(conn net.Conn) {
	//removing connection...
}
