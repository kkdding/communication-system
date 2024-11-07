package main

import (
	"fmt"
	"net"
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
	u.server.BroadCast(u, msg)
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
