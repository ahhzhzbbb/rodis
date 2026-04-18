package main

import (
	_ "net/http/pprof"
	"rodis/internal/server"
)

func main() {
	s := server.NewServer(server.Config{BatchSize: 8})
	s.Start()
}
