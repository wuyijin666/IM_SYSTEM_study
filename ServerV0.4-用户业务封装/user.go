package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string // C 表示是否有数据 回写给当前的客户端
	conn net.Conn
}

// 创建一个用户的API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr, 
		Addr: userAddr, 
		C: make(chan string),
		conn: conn,
	}

	// 启动 监听当前user channel消息的go程
	go user.ListenMessage()
	return user
}

//监听当前 User channel的方法，一旦有消息，就发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
