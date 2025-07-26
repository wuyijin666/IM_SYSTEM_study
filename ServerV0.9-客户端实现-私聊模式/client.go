package main

import (
	"flag"
	"fmt"
	"net"
	"io"
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


// 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

//处理server回执的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	//一旦client.conn中有数据，就直接copy到stdout上，永久阻塞监听
	//io.Copy(os.Stdout, client.conn)
	// 类似下列伪码
	// for {
	// 	buf := make([]byte, 1024)
	// 	client.conn.Read(buf)
	// 	fmt.Println()
	// }
	buf := make([]byte, 4096)
	for {
		n, err := client.conn.Read(buf)
		if n == 0 {
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("读取服务器数据错误:", err)
			return
		}
		// 输出到标准输出
		fmt.Println(string(buf[:n]))

	}
}

func (client *Client) PublicChat() {
	var chatMsg string

	for {
		fmt.Println(">>>>请输入聊天内容，exit退出.")
	    fmt.Scanln(&chatMsg)

		if chatMsg == "exit" {
			break
		}
		// 发送给服务端
		// 消息不为空则发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err: ", err)
				break
			}
		}
	}
}

func (client *Client) Run() {
    for client.Flag != 0 {
		for client.Menu() != true {

		}
		//根据不同的业务选择不同的模式
		switch client.Flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			fmt.Println("2. 私聊模式......")
			break
	    case 3:
		    client.UpdateName()
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

	// 单开一个go程，去处理server的回执消息
	go client.DealResponse()

	// 启动客户端业务
	client.Run()
}
