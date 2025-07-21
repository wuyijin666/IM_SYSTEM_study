package main

import (
	"fmt"
	"net"

)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("dial error: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("连接建立成功")
}
