package ziface

// 定义一个服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// add Router
	AddRouter(id uint32, router IRouter)
	// get message request handler
	GetMsgHandler() IMsgHandler
	// get connection manager handler
	GetConnManager() IConnManager
	// 设置服务器启动钩子函数
	SetHookOnConnStart(func(connection IConnection))
	// 设置服务器停止钩子函数
	SetHookOnConnStop(func(connection IConnection))
	// 调用服务器启动钩子函数
	InvokeHookOnConnStart(connection IConnection)
	// 调用服务器停止钩子函数
	InvokeHookOnConnStop(connection IConnection)
}
