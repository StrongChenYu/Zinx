package ziface

type IDataPack interface {
	// 获取长度
	GetHeadLen() uint32

	// message -> []byte
	Pack(message IMessage) ([]byte, error)

	// []byte -> message
	Unpack([]byte) (IMessage, error)
}
