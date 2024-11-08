package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	mode       int //当前客户端模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		mode:       -1,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client
}

// DealResponse 处理server返回的消息
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，直接copy到stdout标准输出上，永久阻塞监听
	_, err := io.Copy(os.Stdout, client.conn)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}
}

func (client *Client) menu() bool {
	var mode int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&mode)
	if err != nil {
		fmt.Println("Scan error:", err)
		return false
	}

	if mode >= 0 && mode <= 3 {
		client.mode = mode
		return true
	} else {
		fmt.Println(">>>>>>请输入合法范围内的数字<<<<<<")
		return false
	}
}

func (client *Client) UpdateUserName() bool {
	fmt.Println(">>>>>>请输入用户名:")
	_, err := fmt.Scanln(&client.Name)
	if err != nil {
		fmt.Println("Scan error:", err)
		return false
	}

	sendMsg := "rename|" + client.Name + "\n"
	_, err = client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("Write error:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.mode != 0 {
		for !client.menu() {
		}

		switch client.mode {
		case 1:
			// 公聊模式
			fmt.Println("公聊模式...")
			break
		case 2:
			// 私聊模式
			fmt.Println("私聊模式...")
			break
		case 3:
			// 更新用户名
			fmt.Println("更新用户名...")
			client.UpdateUserName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", ServerIp, "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", ServerPort, "设置服务器端口(默认是8080)")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>> 链接服务器失败...")
		return
	}

	// 单独开启一个goroutine处理server的回执消息
	go client.DealResponse()

	fmt.Println(">>>>>> 链接服务器成功...")

	client.Run()

	// 启动客户端业务
	//select {}

}
