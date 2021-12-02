package znet

import "zinx/ziface"

type BaseRouter struct{}

func (router *BaseRouter) BeforeHandler(request ziface.IRequest) {}
func (router *BaseRouter) Handler(request ziface.IRequest)       {}
func (router *BaseRouter) AfterHandler(request ziface.IRequest)  {}
