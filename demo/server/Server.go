package main

import "zinx/znet"

func main() {
	server := znet.NewServer("v1.0")
	server.Serve()
}
