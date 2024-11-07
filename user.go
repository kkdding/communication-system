package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	// 用户所属服务
	server *Server
}

// NewUser 创建用户API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// 启动监听
	go user.ListenMessage()

	return user
}

// Online 用户上线业务
func (u *User) Online() {
	// 用户上线，将用户加入到OnlineMap中
	u.server.mapLock.Lock()

	u.server.OnlineMap[u.Name] = u
	fmt.Println("当前用户列表:")
	for userName, userInfo := range u.server.OnlineMap {
		fmt.Printf("userName: %s, userInfo: %v\n", userName, *userInfo)
	}

	u.server.mapLock.Unlock()

	// 广播当前用户上线消息
	u.server.BroadCast(u, "已上线")

}

// Offline 用户下限业务
func (u *User) Offline() {
	// 用户下线，将用户从OnlineMap中删除
	u.server.mapLock.Lock()

	delete(u.server.OnlineMap, u.Name)
	fmt.Println("当前用户列表:")
	for userName, userInfo := range u.server.OnlineMap {
		fmt.Printf("userName: %s, userInfo: %v\n", userName, *userInfo)
	}

	u.server.mapLock.Unlock()

	// 广播当前用户下线消息
	u.server.BroadCast(u, "已下线")

}

// DoMessage 用户处理消息业务
func (u *User) DoMessage(msg string) {
	if msg == WHO {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineUser := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.C <- onlineUser
		}
		u.server.mapLock.Unlock()
	} else if msg == SELF {
		userInfo := "Info:" + "[" + u.Addr + "]" + u.Name + "\n"
		u.C <- userInfo
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.C <- "当前用户名已被使用\n"
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.C <- "更新用户名为:" + u.Name + "\n"
		}
	} else {
		u.server.BroadCast(u, msg)
	}
}

// ListenMessage 监听User中channel的方法，一旦有消息就发送给客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Listen Error:", err)
			return
		}
	}
}
