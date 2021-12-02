package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

//iServer的接口实现
type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	router    ziface.IRouter
}

func (server *Server) AddRouter(router ziface.IRouter) {
	fmt.Println("Add router...")
	server.router = router
}

// 启动服务器
func (server *Server) Start() {

	fmt.Printf("[Start] Server listening at IP: %s, Port %d, is starting\n", server.IP, server.Port)
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

			connection := NewConnection(tcpConn, uint32(cntId), server.router)
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
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8080,
	}
	return s
}
