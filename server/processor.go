package main

import (
	"fmt"
	"io"
	"net"
	"server/common/message"
	"server/process"
	"server/utils"
)

// 先创建一个结构体

type Processor struct {
	Conn net.Conn
}

// 编写一个serverProcessMes函数
// 功能：根据客户端发送消息的种类不同，决定调用哪个函数来处理
func (this *Processor)serverProcessMes(mes *message.Message) (err error) {
	// 测试客户端是否能接收到群发消息
	//fmt.Println("mes=", mes)

	switch mes.Type {
	case message.LoginMesType :
		// 处理登陆逻辑
		// 创建一个UserProcess实例
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType :
		// 处理注册
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType :
		// 创建一个SmsProcess实例，处理群发消息
		smsProcess := &process.SmsProcess{}
		smsProcess.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

func (this *Processor) process2() (err error) {
	// 循环读取客户端发送的信息
	for {
		// 这里将读取数据包封装成一个readPkg()函数，返回Message，Err
		// 创建一个Transfer实例，完成读包
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器也正常退出...")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		// fmt.Println("mes=", mes)
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}

