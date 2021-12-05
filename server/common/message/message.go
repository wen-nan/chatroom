package message

import "server/model"

const (
	LoginMesType = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType = "SmsMes"
)

// 定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string // 消息类型
	Data string // 消息的数据
}

// 定义两个消息。。 后面需要再增加

type LoginMes struct {
	UserId int // 用户id
	UserPwd string // 用户密码
	UserName string  // 用户名
}

type LoginResMes struct {
	Code int // 返回状态码 500表示该用户还没注册 200表示登陆成功
	// 增加一个字段 保存用户id的切片
	UsersId []int
	Error string //返回错误信息
}

type RegisterMes struct {
	User model.User `json:"user"` // 类型就是User结构体
}

type RegisterResMes struct {
	Code int `json:"code"` // 返回状态码 400表示该用户已经占用 200表示注册成功
	Error string `json:"error"` //返回错误信息
}

// 为了配合服务器端推送用户状态变化消息

type NotifyUserStatusMes struct {
	UserId int `json:"userId"` // 用户id
	Status int `json:"status"` // 用户状态
}

// 增加一个SmsMes，发送的消息

type SmsMes struct {
	Content string `json:"content"` // 内容
	model.User // 匿名结构体，继承
}