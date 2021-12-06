package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type OwnRouter struct {
	ziface.IRouter
}

func (router *OwnRouter) BeforeHandler(request ziface.IRequest) {
	_, err := request.GetConnection().GetTcpConn().Write([]byte("pre handler\n"))
	if err != nil {
		fmt.Println("error occur while invoking before handler")
	}
}
func (router *OwnRouter) Handler(request ziface.IRequest) {
	_, err := request.GetConnection().GetTcpConn().Write([]byte("handler\n"))
	if err != nil {
		fmt.Println("error occur while invoking handler")
	}
}
func (router *OwnRouter) AfterHandler(request ziface.IRequest) {
	_, err := request.GetConnection().GetTcpConn().Write([]byte("after handler\n"))
	if err != nil {
		fmt.Println("error occur while invoking after handler")
	}
}

func main() {
	server := znet.NewServer("v1.0")
	server.AddRouter(&OwnRouter{})
	server.Serve()
}
