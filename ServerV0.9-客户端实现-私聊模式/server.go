package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 全局用户在线列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	// 消息广播的channel
    Message chan string
}

// 创建一个server的接口(对象)
func NewServer(ip string, port int) *Server {
	Server := &Server{
		Ip:   ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}

	return Server
}

// 广播当前上线的用户
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

// 监听当前Message的channel消息的go程 ，目的：一旦Message中有数据，就发送给全部在线的user
func (this *Server) ListenMessger() {
	for {
		msg := <-this.Message

		// 将msg发送给全部在线user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}

		this.mapLock.Unlock()
	}
}

// 处理链接的业务
func (this *Server) Handler(conn net.Conn) {
	// 当前建立的链接
	//fmt.Println("链接建立成功")
	user := NewUser(conn, this)

	// 用户上线，将用户保存到OnlineMap中
	user.OnLine()
    
	// 广播当前上线的用户
	fmt.Printf("[Server] 用户 %s 上线了\n", user.Name)


	// 监听用户是否活跃的channel
    isAlive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		reader := bufio.NewReader(conn)

		for {
			msg , err := reader.ReadString('\n')
			if err == io.EOF {
				// 用户下线
				user.OffLine()
				return
			}

			if err != nil {
				fmt.Println("Read err:", err)
				return 
			}

			// 提取用户消息 字节转字符 （去掉最后一位\n)
			msg = strings.TrimSuffix(msg, "\n")

			if msg == "" {
				continue
			}
			// 广播用户发送的消息
			user.DoMessage(msg)

			// 表示当前用户是一个活跃的
			isAlive <- true
		   }
		}()

	// 当前Handler阻塞
	for {
		select{
		case <-isAlive: 
			// 表示当前用户是活跃的，应该重置下面的定时器
			// 不做任何事，只为了激活select，更新下面的定时器

		case <-time.After(time.Second * 300):
			// 表示用户未活跃 将强制下线该用户
			user.sendMsg("用户：" + user.Name + "被强制下线了")

			//销毁用户资源
			close(user.C)	

			// 关闭连接管道
			conn.Close()
			// 退出当前Handler
			return //runtime.Goexit()

		}
	}
	
}



// 启动服务器的接口
func (this *Server) Start() {
	fmt.Printf("启动服务器，监听地址：%s:%d\n", this.Ip, this.Port)

	// socket listen
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listen error: ", err)
		return 
	}
	// close listen socket
	defer listener.Close()

	// 启动监听Message的go程 因为当server启动，Message就应该是一直被监听的状态
	go this.ListenMessger()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error : ", err)
			continue
		}
		// do handler
		go this.Handler(conn)
	}
	
}
