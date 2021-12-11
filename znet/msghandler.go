package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	router     map[uint32]ziface.IRouter
	WorkerSize int
	Worker     []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	msgHandler := &MsgHandler{
		router:     make(map[uint32]ziface.IRouter),
		WorkerSize: int(utils.GlobalObject.MaxWorkerSize),
		Worker:     make([]chan ziface.IRequest, utils.GlobalObject.MaxWorkerSize),
	}
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

func (msgHandler *MsgHandler) StartWorkerPool() {
	for i := 0; i < msgHandler.WorkerSize; i++ {
		msgHandler.Worker[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerQueueSize)
		go msgHandler.StartWorker(i, msgHandler.Worker[i])
	}
}

func (msgHandler *MsgHandler) StartWorker(id int, queue chan ziface.IRequest) {
	fmt.Printf("Start Worker id: %d\n", id)
	for {
		select {
		case request := <-queue:
			msgHandler.DoMsgHandler(request)
		}
	}
}

func (msgHandler *MsgHandler) HandleRequest(request ziface.IRequest) {
	// hash select algorithm
	selectId := int(request.GetMsgId()) % msgHandler.WorkerSize
	// request send to worker queue
	msgHandler.Worker[selectId] <- request
}

func (msgHandler *MsgHandler) GetWorkerQueue(id int) chan ziface.IRequest {
	return msgHandler.Worker[id]
}

func (msgHandler *MsgHandler) GetWorkerSize() int {
	return msgHandler.WorkerSize
}
