package main

import "rodis/internal/server"

func main() {
	s := server.NewServer(server.Config{})
	s.Start()
}
