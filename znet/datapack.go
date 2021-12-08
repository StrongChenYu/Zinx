package znet

import (
	bytes "bytes"
	binary "encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取长度
func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

// message -> []byte
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if err := binary.Write(buf, binary.LittleEndian, msg.GetLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, msg.GetId()); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// []byte -> message
func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	reader := bytes.NewReader(data)

	msg := &Message{}

	if err := binary.Read(reader, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPacketSize > 0 && msg.Len > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New(" packet exceed max packet size set!")
	}

	return msg, nil
}
