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
	Flag       int
}

func NewClient(serverIp string, serverPort int) *Client { 
	// 创建客户端对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		Flag: 999,
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


//Menu函数 获取用户的输入模式
func (client *Client) Menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出系统")

	fmt.Scanln(&flag)  //从标准输入（键盘）读取用户输入的一行数据

	if flag < 0 || flag > 3 {
		fmt.Println("请输入合法范围内的数字...")
		return false
	}else {
		client.Flag = flag
		return true
	}
}

func (client *Client) Run() {
    for client.Flag != 0 {
		for client.Menu() != true {

		}
		//根据不同的业务选择不同的模式
		switch client.Flag {
		case 1:
			fmt.Println("1. 群聊模式......")
			break
		case 2:
			fmt.Println("2. 私聊模式......")
			break
	    case 3:
		    fmt.Println("3. 更新用户名......")
			break
		}
	}
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
	client.Run()
}
