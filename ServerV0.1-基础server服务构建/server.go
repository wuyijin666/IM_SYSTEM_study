package main

import (
	"fmt"
	"net"
	"sync"
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
	sendMsg := "【" + user.Addr + "]" + user.Name + ":" + msg

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
	user := NewUser(conn)

	// 用户上线，将用户保存到OnlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
    
	// 广播当前上线的用户
	this.BroadCast(user, "已上线")

	fmt.Printf("[Server] 用户 %s 已上线\n", user.Name)

	// 当前handler阻塞
	select {}
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
