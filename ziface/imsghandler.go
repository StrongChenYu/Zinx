package ziface

type IMsgHandler interface {
	DoMsgHandler(IRequest)
	RegisterRouter(uint32, IRouter)
}
