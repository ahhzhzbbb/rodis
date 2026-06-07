package server

import "fmt"

func (s *Server) Banner() {
	fmt.Printf("Welcome to Rodis! Server is running on port %s\n", s.Port)
	fmt.Println(` ____   ___  ____ ___ ____  
|  _ \ / _ \|  _ \_ _/ ___| 
| |_) | | | | | | | |\___ \ 
|  _ <| |_| | |_| | | ___) |
|_| \_\\___/|____/___|____/ `)
	fmt.Println("================={hoangmp}==")

	fmt.Println("Rodis is a Redis-compatible in-memory data structure store written in Go.")
	fmt.Println("Type 'help' for more information about commands.")
	fmt.Println("Email: ahhzhzbbb@gmail.com")
	fmt.Println("Enjoy using Rodis!")
	fmt.Println("..........................................................................")
}
