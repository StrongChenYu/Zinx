package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type OwnRouter1 struct {
	ziface.IRouter
}

type OwnRouter2 struct {
	ziface.IRouter
}

func (router *OwnRouter1) BeforeHandler(request ziface.IRequest) {}
func (router *OwnRouter1) AfterHandler(request ziface.IRequest)  {}
func (router *OwnRouter1) Handler(request ziface.IRequest) {
	err := request.GetConnection().Send(request.GetData(), 0)
	if err != nil {
		fmt.Println("error occur while invoking handler")
	}
}

func (router *OwnRouter2) BeforeHandler(request ziface.IRequest) {}
func (router *OwnRouter2) AfterHandler(request ziface.IRequest)  {}
func (router *OwnRouter2) Handler(request ziface.IRequest) {
	err := request.GetConnection().Send(request.GetData(), 0)
	if err != nil {
		fmt.Println("error occur while invoking handler")
	}
}

func main() {
	server := znet.NewServer("v1.0")
	server.AddRouter(0, &OwnRouter1{})
	server.AddRouter(1, &OwnRouter2{})
	server.Serve()
}
