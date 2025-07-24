package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp string
	ServerPort int
	conn       net.Conn
	Name       string
}

func NewClient(serverIp string, serverPort int) *Client { 
	// 创建客户端对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client

}

var (
	serverIp string
	serverPort int
)

// init函数初始化命令行参数
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip:127.0.0.1")
	flag.IntVar(&serverPort, "port",  8888, "设置服务器端口号:8888")
}

func main(){
	// 命令行解析
	flag.Parse()

	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>> 连接服务器失败")
		return
	} 

	fmt.Println(">>>>连接服务器成功")

	// 启动客户端业务
	select {
		
	}
}
