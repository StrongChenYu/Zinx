package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer
	Host      string
	Port      uint32
	Name      string
	Version   string

	MaxPacketSize      uint32
	MaxConn            uint32
	MaxWorkerSize      uint32
	MaxWorkerQueueSize uint32
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		TcpServer:          nil,
		Host:               "0.0.0.0",
		Port:               8080,
		Name:               "Zinx Server",
		Version:            "0",
		MaxPacketSize:      4096,
		MaxConn:            1000,
		MaxWorkerSize:      10,
		MaxWorkerQueueSize: 1024,
	}

	GlobalObject.Reload()
}

func (gb *GlobalObj) Reload() {
	// 读文件
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	// 将配置文件解析到对象中
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}
