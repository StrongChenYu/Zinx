package ziface

type IMessage interface {
	GetLen() uint32
	GetId() uint32
	GetData() []byte

	SetLen(uint32)
	SetId(uint32)
	SetData([]byte)
}
