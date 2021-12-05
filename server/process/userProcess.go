package process

import (
	"encoding/json"
	"fmt"
	"net"
	"server/common/message"
	"server/model"
	"server/utils"
)

type UserProcess struct {
	Conn net.Conn
	// 增加一个字段，表示该Conn是哪个用户
	UserId int
}

// 通知所有在线用户的方法, userId要通知其他人自己上线了

func (this *UserProcess)NotifyOthersOnlineUser(userId int) {
	// 遍历onlineUsers,然后一个一个的发送 NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		// 过滤自己
		if id == userId {
			continue
		}
		// 开始通知 [单独写一个方法]
		up.NotifyMeOnline(userId)
	}

}

func (this *UserProcess)NotifyMeOnline(userId int) {
	// 组装我们的消息NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	// 将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 将序列化后的mes赋值给data
	mes.Data = string(data)
	// 再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 发送，创建transfer实例，发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

func (this *UserProcess)ServerProcessRegister(mes *message.Message) (err error) {
	// 先从mes中取出mes.data,并直接反序列化为registerMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	// 先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	// 再声明一个registerResMes并完成赋值
	var registerResMes message.RegisterResMes

	// 现在需要到Redis数据库完成注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 400
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}
	// 将loginResMes序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 将data赋值给resMes
	resMes.Data = string(data)
	// 对resMes进行序列化,准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 发送data,将其封装到writePkg函数
	// 因为使用了分层模式，先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

// 编写serverProcessLogin 函数，处理登陆请求

func (this *UserProcess)ServerProcessLogin(mes *message.Message) (err error) {
	// 核心代码
	// 先从mes中取出mes.data,并直接反序列化为loginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	// 先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	// 再声明一个loginResMes并完成赋值
	var loginResMes message.LoginResMes

	// 现在需要到Redis数据库完成验证
	// 使用model.MyUserDao去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}
	} else {
		loginResMes.Code = 200 // 表示合法
		// 这里用户登陆成功，把该登陆成功的用户放入到userMgr中
		// 将登陆成功的用户id赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		// 通知其他在线用户，自己上线
		this.NotifyOthersOnlineUser(loginMes.UserId)
		// 将当前在线用户的id放入到loginResMes.UsersId
		// 遍历userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登陆成功")
	}
	// 如果用户id为100， 密码为123456 就认为合法
	//if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	//	// 合法
	//	loginResMes.Code = 200 // 表示合法
	//} else {
	//	// 不合法
	//	loginResMes.Code = 500 // 500表示该用户不存在
	//	loginResMes.Error = "该用户不存在，请注册再使用..."
	//}

	// 将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 将data赋值给resMes
	resMes.Data = string(data)
	// 对resMes进行序列化,准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 发送data,将其封装到writePkg函数
	// 因为使用了分层模式，先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
