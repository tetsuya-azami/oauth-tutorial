package main

import "oauth-tutorial/internal/server"

func main() {
	s := server.NewServer()
	s.Start()
}
