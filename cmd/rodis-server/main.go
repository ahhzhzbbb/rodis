package main

import (
	_ "net/http/pprof"
	"rodis/internal/server"
)

func main() {
	s := server.NewServer(server.Config{BatchSize: 64})
	s.Start()
}
