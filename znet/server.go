package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//iServer的接口实现
type Server struct {
	Name       string
	IPVersion  string
	IP         string
	Port       int
	MsgHandler ziface.IMsgHandler
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

		cntId := 0
		for {
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept tcp client error: ", err)
				continue
			}

			connection := NewConnection(tcpConn, uint32(cntId), server.MsgHandler)
			cntId++

			go connection.Start()
		}
	}()
}

// 停止服务器
func (server *Server) Stop() {

}

// 运行服务器
func (server *Server) Serve() {
	server.Start()
	select {}
}

func NewServer(name string) ziface.IServer {
	var s = &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       int(utils.GlobalObject.Port),
		MsgHandler: NewMsgHandler(),
	}
	return s
}
