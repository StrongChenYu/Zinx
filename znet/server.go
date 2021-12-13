package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//iServer的接口实现
type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        int
	MsgHandler  ziface.IMsgHandler
	ConnManager ziface.IConnManager
	// 启动钩子函数
	OnStart func(connection ziface.IConnection)
	// 停止钩子函数
	OnStop func(connection ziface.IConnection)
}

func (server *Server) AddRouter(id uint32, router ziface.IRouter) {
	fmt.Println("Add router...")
	server.MsgHandler.RegisterRouter(id, router)
}

// 启动服务器
func (server *Server) Start() {
	fmt.Printf("[Start] Server name: %s listening at IP: %s, Port %d, is starting\n", server.Name, server.IP, server.Port)
	fmt.Printf("[%s] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	go func() {
		addr, err := net.ResolveTCPAddr(server.IPVersion, fmt.Sprintf("%s:%d", server.IP, server.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error")
			return
		}

		listener, err := net.ListenTCP(server.IPVersion, addr)
		if err != nil {
			fmt.Printf("Listen at address: %s:%d error!", server.IP, server.Port)
			return
		}

		// 连接处理goroutine
		server.MsgHandler.StartWorkerPool()

		cntId := 0
		for {
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept tcp client error: ", err)
				continue
			}

			if server.ConnManager.Len() >= int(utils.GlobalObject.MaxConn) {
				err := tcpConn.Close()
				if err != nil {
					fmt.Println("Accept tcp client error: ", err)
					continue
				}
				continue
			}

			connection := NewConnection(tcpConn, uint32(cntId), server)
			server.ConnManager.Add(connection)

			cntId++

			go connection.Start()
		}
	}()
}

// 停止服务器
func (server *Server) Stop() {
	server.ConnManager.ClearAll()
}

// 运行服务器
func (server *Server) Serve() {
	server.Start()
	select {}
}

func (server *Server) GetMsgHandler() ziface.IMsgHandler {
	return server.MsgHandler
}
func (server *Server) GetConnManager() ziface.IConnManager {
	return server.ConnManager
}

func (server *Server) SetHookOnConnStart(f func(connection ziface.IConnection)) {
	server.OnStart = f
}

// 设置服务器停止钩子函数
func (server *Server) SetHookOnConnStop(f func(connection ziface.IConnection)) {
	server.OnStop = f
}

// 调用服务器启动钩子函数
func (server *Server) InvokeHookOnConnStart(connection ziface.IConnection) {
	if server.OnStart != nil {
		fmt.Println("Invoke Server Start hook!")
		server.OnStart(connection)
	}
}

// 调用服务器停止钩子函数
func (server *Server) InvokeHookOnConnStop(connection ziface.IConnection) {
	if server.OnStop != nil {
		fmt.Println("Invoke Server Start hook!")
		server.OnStop(connection)
	}
}

func NewServer() ziface.IServer {
	var s = &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        int(utils.GlobalObject.Port),
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}
