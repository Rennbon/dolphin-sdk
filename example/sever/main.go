package main

import (
	"github.com/2se/dolphin-sdk/mock"
	"github.com/2se/dolphin-sdk/server"
	"time"
)

func main() {
	c := &server.Config{
		AppName:         "userApp",
		Address:         "192.168.10.169:8848",
		WriteBufSize:    32 * 1024,
		ReadBufSize:     32 * 1024,
		ConnTimeout:     time.Second * 10,
		DolphinHttpAddr: "http://192.168.9.130:9527",
		DolphinGrpcAddr: "192.168.9.130:9528",
		RequestTimeout:  time.Second * 30,
	}

	//启动并注册到dolphin
	//1. 启动dolphin
	//2. 启动server
	//3. 启动client
	//server.Start(c, mock.MkService)
	//只启动grpc
	//1. 启动server
	//2. 启动client
	server.StartGrpcOnly(c, mock.MkService)
}
