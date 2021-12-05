package model

import "net"

type CurUser struct {
	Conn net.Conn
	User // 匿名结构体
}


