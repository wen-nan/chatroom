package main

import (
	"fmt"
	"net"
	"server/model"
	"time"
)

// 处理和客户端的通讯
func mainProcess(conn net.Conn) {
	// 延时关闭
	defer conn.Close()
	// 调用总控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器端通讯的协程错误 err=", err)
		return
	}
}

// 编写一个函数，完成对userDao初始化
func initUserDao() {
	// 这里的pool 就是一个全局变量
	// 但需要注意初始化顺序问题
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	// 当服务器启动时，就初始化redis连接池
	initPool("localhost:6379", 16, 0, 300 * time.Second)
	initUserDao()
	// 提示信息
	fmt.Println("服务器[新结构]在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	// 一旦监听成功，就等待客户端连接服务器
	for {
		fmt.Println("等待客户端连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		// 连接成功，则启动一个协程和客户端保持通讯...
		go mainProcess(conn)
	}
}