package znet

// head(4 byte) | Id(4 byte) | Data
type Message struct {
	Len  uint32
	Id   uint32
	Data []byte
}

func (msg *Message) GetLen() uint32 {
	return msg.Len
}
func (msg *Message) GetId() uint32 {
	return msg.Id
}
func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetLen(len uint32) {
	msg.Len = len
}
func (msg *Message) SetId(id uint32) {
	msg.Id = id
}
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
