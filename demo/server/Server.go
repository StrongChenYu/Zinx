package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type OwnRouter struct {
	ziface.IRouter
}

func (router *OwnRouter) BeforeHandler(request ziface.IRequest) {}
func (router *OwnRouter) AfterHandler(request ziface.IRequest)  {}

func (router *OwnRouter) Handler(request ziface.IRequest) {
	err := request.GetConnection().Send(request.GetData(), 0)
	if err != nil {
		fmt.Println("error occur while invoking handler")
	}
}

func main() {
	server := znet.NewServer("v1.0")
	server.AddRouter(&OwnRouter{})
	server.Serve()
}
