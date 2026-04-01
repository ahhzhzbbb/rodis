package main

import "rodis/internal/server"

func main() {
	s := server.NewServer(server.Config{Port: ":6379"})
	s.Start()
}
