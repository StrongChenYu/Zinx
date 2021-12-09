package znet

import (
	"fmt"
	"zinx/ziface"
)

type MsgHandler struct {
	router map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	msgHandler := &MsgHandler{router: make(map[uint32]ziface.IRouter)}
	return msgHandler
}

func (msgHandler *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	router, ok := msgHandler.router[request.GetMsgId()]
	if !ok {
		fmt.Printf("router corresponding to %d does not exist.\n", request.GetMsgId())
		return
	}
	router.BeforeHandler(request)
	router.Handler(request)
	router.AfterHandler(request)
}

func (msgHandler *MsgHandler) RegisterRouter(id uint32, router ziface.IRouter) {
	if _, ok := msgHandler.router[id]; ok {
		panic("router already exist!")
	}
	msgHandler.router[id] = router
}
