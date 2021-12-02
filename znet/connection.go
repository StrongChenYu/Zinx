package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	// 当前连接的套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnId uint32
	// 当前的连接状态
	isClosed bool
	// 当前连接所绑定的Router
	router ziface.IRouter
	// 告知当前链接已经退出/停止 channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	s := Connection{
		Conn:     conn,
		ConnId:   connId,
		isClosed: false,
		router:   router,
		ExitChan: make(chan bool, 1),
	}
	return &s
}

// 读业务
func (conn *Connection) StartReader() {
	fmt.Printf("Connection %d Reader goroutine is starting...\n", conn.ConnId)
	defer fmt.Printf("Connection %d Reader is exited...\n", conn.ConnId)
	defer conn.Stop()

	for {
		buf := make([]byte, 512)
		_, err := conn.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error", err)
			continue
		}

		request := Request{
			Conn: conn,
			data: buf,
		}

		go func(request ziface.IRequest) {
			conn.router.BeforeHandler(request)
			conn.router.Handler(request)
			conn.router.AfterHandler(request)
		}(&request)

	}
}

// 写业务
func (conn *Connection) StartWriter() {
	// TODO: 读写分离
	return
}

//开启链接
func (conn *Connection) Start() {
	fmt.Println("Conn Start(), id is: ", conn.ConnId)
	// 读写分离
	go conn.StartReader()
	// go conn.StartWriter()
}

// 停止连接
func (conn *Connection) Stop() {
	fmt.Println("Connection stop(), id is:", conn.ConnId)

	if conn.isClosed {
		return
	}

	conn.isClosed = true
	if err := conn.Conn.Close(); err != nil {
		fmt.Println("Connection stop error, id is:", conn.ConnId)
	}
	close(conn.ExitChan)
}

// 获取对应的tcp 连接
func (conn *Connection) GetTcpConn() *net.TCPConn {
	return conn.Conn
}

// 获取这个connection的id
func (conn *Connection) GetConnID() uint32 {
	return conn.ConnId
}

// 连接的远程id
func (conn *Connection) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

// 发送数据是否成功
func (conn *Connection) Send(data []byte) error {
	return nil
}
