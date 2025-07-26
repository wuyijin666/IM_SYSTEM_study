package main

import (
	"fmt"
	"net"
	"strings"
	"time"
	"sync"
)

type User struct {
	Name string
	Addr string
	C    chan string // C 表示是否有数据 回写给当前的客户端
	conn net.Conn
	connMutex sync.Mutex  // 添加连接写入锁

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

//给当前user对应的客户端发送对应的消息
func (this *User) sendMsg(msg string) {
	this.connMutex.Lock()
	defer this.connMutex.Unlock()
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n" 
	}
	this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string){
	trimmedMsg := strings.TrimSpace(msg)

	if trimmedMsg == "who"{
		this.server.mapLock.RLock()
		// 查询当前用户都有哪些
		onlineMsg := "当前在线用户：\n"
		for _, user := range this.server.OnlineMap {
            onlineMsg += " " + user.Name + "\n"
			
		}
        
		this.server.mapLock.RUnlock()
		this.sendMsg(onlineMsg)
	}else if len(trimmedMsg) > 7 && trimmedMsg[:7] == "rename|" {
		// 消息格式 "rename|kolar"
		// 直接取 | 后面的部分 
		newName := trimmedMsg[7:]
		// 判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
              this.sendMsg("当前用户名已被使用")
		}else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.Name = newName
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.sendMsg("您已修改用户名" + newName + "\n")
		}
	}else if len(trimmedMsg) > 3 && trimmedMsg[:3] == "to|" { 
		// 消息格式 "to|张三|消息内容"
		// 1. 获取对方用户名
		remoteName := strings.Split(trimmedMsg, "|")[1]
		if remoteName == "" {
			this.sendMsg("消息格式不正确，请使用\"to|张三|消息内容\"格式。\n")
			return 
		}
		// 2. 根据用户名，得到对方的User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.sendMsg("该用户不存在, 请重新输入")
			return 
		}

		// 3. 获取消息内容，通过对方的User对象，将消息内容发送出去
		content := strings.Split(trimmedMsg, "|")[2]
		if content == "" {
			this.sendMsg("发送新的消息内容不能为空")
			return 
		}
		remoteUser.sendMsg(this.Name + "对你说：" + content)

	}else {
		this.server.BroadCast(this, msg)
    }
}
   
//监听当前 User channel的方法，一旦有消息，就发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		// this.conn.Write([]byte(msg + "\n"))
		this.sendMsg(msg)
	}
}
