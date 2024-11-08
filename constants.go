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
