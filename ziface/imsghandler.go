package ziface

type IMsgHandler interface {
	DoMsgHandler(IRequest)
	RegisterRouter(uint32, IRouter)
	StartWorkerPool()
	GetWorkerQueue(id int) chan IRequest
	GetWorkerSize() int
	HandleRequest(IRequest)
}
