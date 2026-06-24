package server

import (
	"fmt"
	"io"
	"net"
	"rodis/internal/command"
	"rodis/internal/engine"
	"rodis/internal/factory"
	"rodis/internal/protocol/resp"
	"strings"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	aof, err := engine.NewAof("../appendonly.aof")
	if err != nil {
		fmt.Println("failed to create aof")
		return
	}
	// info, err := os.Stat("../appendonly.aof")
	// fmt.Println(info, err)

	rp := resp.NewResp(conn)

	aofResp := resp.NewResp(aof.File)
	go func() {
		request, err := aofResp.ParseRESP()
		if err != nil {
			return
		}
		s.handleRequest(request)
	}()

	for {
		count := 0

		for count < s.BatchSize {
			request, err := rp.ParseRESP()
			if err != nil {
				// fmt.Printf("READ ERROR: %T | %v\n", err, err)
				if err == io.EOF {
					fmt.Printf("client %s disconnected\n", conn.RemoteAddr())
					return
				}
				// s.removeConnection(conn)
				return
			}
			fmt.Printf("request: %v\n", request)

			if len(request.Array) == 0 {
				continue
			}

			typeOfCommand := strings.ToUpper(request.Array[0].Bulk)
			if typeOfCommand == "SET" ||
				typeOfCommand == "DEL" ||
				typeOfCommand == "INCR" ||
				typeOfCommand == "APPEND" ||
				typeOfCommand == "RPUSH" ||
				typeOfCommand == "LPUSH" ||
				typeOfCommand == "RPOP" ||
				typeOfCommand == "LPOP" ||
				typeOfCommand == "LINSERT" {
				aof.Write(request, *aofResp)
			}

			response := s.handleRequest(request)
			fmt.Printf("response: %v\n", response)

			if err := rp.Marshal(response); err != nil {
				return
			}

			count++
			if !rp.HasBufferedData() {
				break
			}
		}
		if err := rp.FlushWriter(); err != nil { //flush
			return
		}
	}
}

func (s *Server) handleRequest(request resp.Payload) resp.Payload {
	var result resp.Payload
	typeOfCommand := strings.ToUpper(request.Array[0].Bulk)

	comCreator, ok := factory.CommandRegistry[typeOfCommand]
	if !ok {
		result = resp.NewError(fmt.Sprintf("ERR unknown command '%s', with args beginning with: ", typeOfCommand))
	} else {
		comm := comCreator()
		result = comm.Execute(request.Array[1:], command.NewCommandContext(s.kv))
	}
	return result
}

func (s *Server) removeConnection(conn net.Conn) {
	//removing connection from registry...
}
