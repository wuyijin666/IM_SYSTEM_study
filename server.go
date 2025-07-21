package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	Server := &Server{
		Ip:   ip,
		Port: port,
	}
	return Server
}

func (this *Server) Handler(conn net.Conn) {
	// 当前建立的链接
	fmt.Println("链接建立成功")
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
	defer listener.Close()

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
	
	// close listen socket
}
