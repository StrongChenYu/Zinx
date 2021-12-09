package znet

import "zinx/ziface"

type Request struct {
	Conn ziface.IConnection
	Msg  ziface.IMessage
}

func (request *Request) GetConnection() ziface.IConnection {
	return request.Conn
}

// 这个request中的data
func (request *Request) GetData() []byte {
	return request.Msg.GetData()
}

func (request *Request) GetMsgId() uint32 {
	return request.Msg.GetId()
}
