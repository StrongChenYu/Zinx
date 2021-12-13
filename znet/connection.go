package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
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
	Server ziface.IServer
	// 告知当前链接已经退出/停止 channel
	ExitChan chan bool
	// 发送的信息通道，只需要往这个通道里面写就可以了
	WriteChan chan []byte
	// 发送消息的缓冲
	WriteBufChan chan []byte
	// 默认的拆包器
	DataPack ziface.IDataPack
	// 属性
	properties map[string]interface{}
	// 属性锁
	pMutex sync.RWMutex
}

func NewConnection(conn *net.TCPConn, connId uint32, server ziface.IServer) *Connection {
	s := Connection{
		Conn:         conn,
		ConnId:       connId,
		IsClosed:     false,
		Server:       server,
		ExitChan:     make(chan bool),
		WriteChan:    make(chan []byte, 1),
		WriteBufChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		DataPack:     NewDataPack(),
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

		if conn.Server.GetMsgHandler().GetWorkerSize() > 0 {
			// use pool to handle request
			conn.Server.GetMsgHandler().HandleRequest(&request)
		} else {
			// just handle request in one goroutine, after request is processed, goroutine will be destroyed
			go conn.Server.GetMsgHandler().DoMsgHandler(&request)
		}
	}
}

// 写业务
func (conn *Connection) StartWriter() {
	fmt.Printf("Connection %d Writer goroutine is starting...\n", conn.ConnId)
	defer fmt.Printf("Connection %d Writer is exited...\n", conn.ConnId)

	for {
		select {
		case data := <-conn.WriteChan:
			if _, err := conn.Conn.Write(data); err != nil {
				fmt.Println("write error: ", err)
				continue
			}
		case data, ok := <-conn.WriteBufChan:
			if ok {
				if _, err := conn.Conn.Write(data); err != nil {
					fmt.Println("write error: ", err)
					continue
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-conn.ExitChan:
			return
		}
	}
}

// 发送数据是否成功
func (conn *Connection) Send(data []byte, msgId uint32) error {
	if conn.IsClosed {
		fmt.Println("connection already close: ")
		return errors.New("connection already close")
	}

	packet := NewMessagePacket(data, msgId)

	binaryData, err := conn.DataPack.Pack(packet)
	if err != nil {
		fmt.Println("pack message error: ", err)
		return err
	}

	conn.WriteChan <- binaryData
	return nil
}

func (conn *Connection) SendBuf(data []byte, msgId uint32) error {
	if conn.IsClosed {
		fmt.Println("connection already close: ")
		return errors.New("connection already close")
	}

	packet := NewMessagePacket(data, msgId)

	binaryData, err := conn.DataPack.Pack(packet)
	if err != nil {
		fmt.Println("pack message error: ", err)
		return err
	}

	conn.WriteBufChan <- binaryData
	return nil
}

// 增加属性
func (conn *Connection) AddProperty(key string, value interface{}) {
	conn.pMutex.Lock()
	defer conn.pMutex.Unlock()

	conn.properties[key] = value
}

// 获取属性
func (conn *Connection) GetProperty(key string) interface{} {
	conn.pMutex.RLock()
	defer conn.pMutex.RUnlock()

	return conn.properties[key]
}

// 删除属性
func (conn *Connection) DeleteProperty(key string) {
	conn.pMutex.Lock()
	defer conn.pMutex.Unlock()

	delete(conn.properties, key)
}

//开启链接
func (conn *Connection) Start() {
	fmt.Println("Conn Start(), Id is: ", conn.ConnId)
	// 读写分离
	go conn.StartReader()
	go conn.StartWriter()

	conn.Server.InvokeHookOnConnStart(conn)
}

// 停止连接
func (conn *Connection) Stop() {
	fmt.Println("Connection stop(), Id is:", conn.ConnId)

	if conn.IsClosed {
		return
	}

	conn.IsClosed = true
	conn.ExitChan <- true

	conn.Server.InvokeHookOnConnStop(conn)

	// 删除链接管理池中的链接
	conn.Server.GetConnManager().Delete(conn)

	if err := conn.Conn.Close(); err != nil {
		panic("Connection stop error")
	}

	close(conn.ExitChan)
	close(conn.WriteChan)
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
