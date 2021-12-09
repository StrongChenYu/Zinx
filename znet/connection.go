package znet

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	// 当前连接的套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnId uint32
	// 当前的连接状态
	IsClosed bool
	// 当前连接管理的router
	MsgHandler ziface.IMsgHandler
	// 告知当前链接已经退出/停止 channel
	ExitChan chan bool
	// 默认的拆包器
	DataPack ziface.IDataPack
}

func NewConnection(conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandler) *Connection {
	s := Connection{
		Conn:       conn,
		ConnId:     connId,
		IsClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		DataPack:   NewDataPack(),
	}
	return &s
}

// 读业务
func (conn *Connection) StartReader() {
	fmt.Printf("Connection %d Reader goroutine is starting...\n", conn.ConnId)
	defer fmt.Printf("Connection %d Reader is exited...\n", conn.ConnId)
	defer conn.Stop()

	for {
		// 先读取头部数据
		headBuf := make([]byte, conn.DataPack.GetHeadLen())
		_, err := io.ReadFull(conn.Conn, headBuf)
		if err != nil {
			fmt.Println("read head error: ", err)
			break
		}

		// 然后解压
		msg, err := conn.DataPack.Unpack(headBuf)
		if err != nil {
			fmt.Println("unpack buf error: ", err)
			break
		}

		// 然后根据长度继续读
		var bodyBuf []byte
		if msg.GetLen() > 0 {
			bodyBuf = make([]byte, msg.GetLen())
			_, err = io.ReadFull(conn.Conn, bodyBuf)
			if err != nil {
				fmt.Println("read body error: ", err)
				break
			}
		}

		// 封装成为request
		msg.SetData(bodyBuf)
		request := Request{
			Conn: conn,
			Msg:  msg,
		}

		go conn.MsgHandler.DoMsgHandler(&request)

	}
}

// 发送数据是否成功
func (conn *Connection) Send(data []byte, msgId uint32) error {
	packet := NewMessagePacket(data, msgId)

	binaryData, err := conn.DataPack.Pack(packet)
	if err != nil {
		fmt.Println("pack message error: ", err)
		return err
	}

	_, err = conn.Conn.Write(binaryData)
	if err != nil {
		fmt.Println("send connection error: ", err)
	}

	return nil
}

// 写业务
func (conn *Connection) StartWriter() {
	// TODO: 读写分离
	return
}

//开启链接
func (conn *Connection) Start() {
	fmt.Println("Conn Start(), Id is: ", conn.ConnId)
	// 读写分离
	go conn.StartReader()
	// go conn.StartWriter()
}

// 停止连接
func (conn *Connection) Stop() {
	fmt.Println("Connection stop(), Id is:", conn.ConnId)

	if conn.IsClosed {
		return
	}

	conn.IsClosed = true
	if err := conn.Conn.Close(); err != nil {
		fmt.Println("Connection stop error, Id is:", conn.ConnId)
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
