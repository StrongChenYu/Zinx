package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("Client start...")

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Connect error", err)
		return
	}

	dp := znet.NewDataPack()
	for {
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

		sendData1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("pack error: ", err)
			continue
		}

		sendData2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("pack error: ", err)
			continue
		}

		sendData1 = append(sendData1, sendData2...)
		_, err = conn.Write(sendData1)
		if err != nil {
			fmt.Println("connection write error: ", err)
			continue
		}

		fmt.Println("Write to server success!")

		headBuf := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headBuf)
		if err != nil {
			fmt.Println("read head error: ", err)
			break
		}

		// 然后解压
		msg, err := dp.Unpack(headBuf)
		if err != nil {
			fmt.Println("unpack buf error: ", err)
			break
		}

		// 然后根据长度继续读
		var bodyBuf []byte
		if msg.GetLen() > 0 {
			bodyBuf = make([]byte, msg.GetLen())
			_, err = io.ReadFull(conn, bodyBuf)
			if err != nil {
				fmt.Println("read body error: ", err)
				break
			}

			fmt.Printf("message from server: %s\n", bodyBuf)
		}

		time.Sleep(1 * time.Second)
	}
}
