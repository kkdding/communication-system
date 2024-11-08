package main

import "time"

// 用户操作相关常量
const (
	OpWho    string = "who"
	OpSelf   string = "self"
	OpRename string = "rename|"
	OpTo     string = "to|"
)

// 时间相关常量
const (
	TimeOut   = time.Second * 10
	TimeDelay = time.Second * 2
)

// 服务器地址和端口常量
const (
	ServerIp   string = "127.0.0.1"
	ServerPort int    = 8080
)
