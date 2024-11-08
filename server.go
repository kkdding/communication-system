package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个 server 的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMessage 监听Message广播消息，发给全部到在线User
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// BroadCast 广播消息
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Name + "] " + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// 链接建立成功
	fmt.Printf("Connect Success!!! Connected Address: %s\n", conn.LocalAddr().String())

	user := NewUser(conn, s)
	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)

			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err == io.EOF {
				fmt.Println("Conn Read Error:", err)
				return
			}

			// 提取用户的消息(去除'\n')
			msg := string(buf[:n-1])

			// 用户针对msg进行消息处理
			user.DoMessage(msg)

			// 用户发送任意消息代表活跃
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
			// 当前用户活跃，重制定时器
			// 激活select

		case <-time.After(time.Second * 10):
			// 用户无操作超时
			// 强制关闭当前User
			user.C <- "长时间未操作，强制下线"
			time.Sleep(time.Second * 2)
			close(user.C)
			err := conn.Close()
			if err != nil {
				fmt.Println("Close Error:", err)
			}
			return
		}
	}
}

// Start 启动 server 的接口
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// close socket listen
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("listener.Close err:", err)
			return
		}
	}(listener)

	// 启动监听Message
	go s.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// do handler
		go s.Handler(conn)
	}

}
