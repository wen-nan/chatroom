package process

import (
	"client/common/message"
	"client/model"
	"client/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
	// 字段...
}

// 注册方法

func (this *UserProcess)Register(userId int,
	userPwd string, userName string) (err error) {
	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dail err=", err)
		return
	}
	// 延时关闭
	defer conn.Close()

	// 准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.RegisterMesType
	// 创建一个RegisterMessage结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 把dada赋给mes.Data字段
	mes.Data = string(data)
	// 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 这里还需要处理服务器返回的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	// 发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}
	// 将mes的data部分反序列化为registerResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功,可以重新登陆")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}


// 关联一个用户登陆的方法
// 完成一个登陆校验

func (this *UserProcess)Login(userId int, userPwd string) (err error) {
	// 开始定协议。。。
	//fmt.Printf("userId=%d userPwd=%s\n", userId, userPwd)
	//
	//return nil

	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dail err=", err)
		return
	}
	// 延时关闭
	defer conn.Close()

	// 准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.LoginMesType
	// 创建一个loginMessage结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	// 将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 把dada赋给mes.Data字段
	mes.Data = string(data)

	// 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 此时data就是我们要发送的消息
	// 1.先把data的长度发送给服务器
	// 先获取data长度->转成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	//fmt.Printf("客户端发送消息长度成功!长度=%d,内容=%s\n", len(data), string(data))

	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	// 这里还需要处理服务器返回的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}
	// 将mes的data部分反序列化为loginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		// 初始化curUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		//fmt.Println("用户登陆成功")
		// 可以显示当前用户在线列表，遍历loginResMes.UsersId
		fmt.Println("当前用户在线列表如下")
		for _, v := range loginResMes.UsersId {
			// 不显示自己
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			// 初始化客户端onlineUsers
			user := &model.User{
				UserId: v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Println()
		// 这里需要再客户端启动一个协程
		// 该协程保持和服务器端的通讯,如果服务器有数据推送，
		// 则接收并显示在客户端终端
		go serverProcessMes(conn)

		// 循环显示登陆成功后的菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

