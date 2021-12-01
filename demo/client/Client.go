package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client start...")

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Connect error", err)
		return
	}
	for {
		fmt.Printf("Write to server: %s\n", "hello world!")

		_ ,err := conn.Write([]byte("hello world!"))

		if err != nil {
			fmt.Println("write to server error: ", err)
			continue
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)

		if err != nil {
			fmt.Println("read from server error: ", err)
			continue
		}

		fmt.Printf("Receive from server: %s\n", buf[:cnt])

		time.Sleep(1 * time.Second)
	}
}
