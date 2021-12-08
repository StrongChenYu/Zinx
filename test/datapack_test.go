package test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"zinx/znet"
)

func TestDataPackServer(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		fmt.Println("test error: ", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err: ", err)
			continue
		}

		go func(conn net.Conn) {
			for {
				dp := znet.NewDataPack()
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)

				if err != nil {
					fmt.Println("read error: ", err)
					return
				}

				msg, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("unpack error: ", err)
					return
				}

				if msg.GetLen() > 0 {
					msg.SetData(make([]byte, msg.GetLen()))
					_, err = io.ReadFull(conn, msg.GetData())
					if err != nil {
						fmt.Println("read error: ", err)
						return
					}
				}
				fmt.Printf("=> Recv Msg: ID=%d, len=%d, data=%s\n", msg.GetId(), msg.GetLen(), string(msg.GetData()))
			}
		}(conn)
	}
}

func TestDataPackTest(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		fmt.Println("connection error: ", err)
		return
	}

	msg1 := &znet.Message{
		Len:  4,
		Id:   5,
		Data: []byte{'h', 'e', 'l', 'l'},
	}

	msg2 := &znet.Message{
		Len:  6,
		Id:   5,
		Data: []byte{'h', 'e', 'l', 'l', 'o', 'o'},
	}

	dp := znet.NewDataPack()

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack error: ", err)
		return
	}

	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack error: ", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)

	_, err = conn.Write(sendData1)

	if err != nil {
		fmt.Println("connection write error: ", err)
		return
	}
}
