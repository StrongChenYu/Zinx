package znet

import "zinx/ziface"

type Request struct {
	Conn ziface.IConnection
	data []byte
}

func (request *Request) GetConnection() ziface.IConnection {
	return request.Conn
}

// 这个request中的data
func (request *Request) GetData() []byte {
	return request.data
}
