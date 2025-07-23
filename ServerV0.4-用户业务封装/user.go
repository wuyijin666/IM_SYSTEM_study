package main

import (
	"fmt"
	"net"
	"time"

)

type User struct {
	Name string
	Addr string
	C    chan string // C 表示是否有数据 回写给当前的客户端
	conn net.Conn

	server *Server    // User类型新增server关联
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr, 
		Addr: userAddr, 
		C: make(chan string),
		conn: conn,

		server: server,
	}

	// 启动 监听当前user channel消息的go程
	go user.ListenMessage()
	return user
}

// 用户上线功能
func (this *User) OnLine() {
    this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "已上线")	
}

// 用户下线功能
func (this *User) OffLine() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "已下线")
	fmt.Printf("[%s] 用户 %s 下线\n", time.Now().Format("2006-01-02 15:04:05"), this.Name)  // 打印用户下线信息
}


// 用户处理消息的业务
func (this *User) DoMessage(msg string){
	this.server.BroadCast(this, msg)
}

//监听当前 User channel的方法，一旦有消息，就发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
