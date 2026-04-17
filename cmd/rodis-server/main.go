package main

import (
	_ "net/http/pprof"
	"rodis/internal/server"
)

func main() {
	s := server.NewServer(server.Config{BatchSize: 16})
	s.Start()
}
