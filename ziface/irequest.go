package ziface

// 这是为了封装一个请求的消息
type IRequest interface {
	// 获取这个消息的链接
	GetConnection() IConnection
	// 这个request中的data
	GetData() []byte
}
