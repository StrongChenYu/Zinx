package ziface

import "net"

type IConnection interface {
	// 开启连接
	Start()
	// 停止连接
	Stop()
	// 获取对应的tcp 连接
	GetTcpConn() *net.TCPConn
	// 获取这个connection的id
	GetConnID() uint32
	// 连接的远程id
	RemoteAddr() net.Addr
	// 发送数据是否成功
	Send(data []byte, msgId uint32) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error
