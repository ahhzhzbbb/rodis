package main

import (
	_ "net/http/pprof"
	"rodis/internal/server"
)

func main() {
	s := server.NewServer(server.DefaultConfig())
	s.Start()
}
