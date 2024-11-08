package main

import "fmt"

func main() {
	fmt.Println("hello")
	server := NewServer(ServerIp, ServerPort)
	server.Start()
}
