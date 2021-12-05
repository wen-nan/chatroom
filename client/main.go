package main

import (
	"client/process"
	"fmt"
	"os"
)

// 定义两个变量，一个表示用户ID,一个表示用户密码
var userId int
var userPwd string
var userName string

func main() {
	// 接收用户的选择
	var key int
	// 判断是否还继续显示菜单
	//var loop = true

	for {
		fmt.Println("----------欢迎登陆多人聊天系统----------")
		fmt.Println("\t\t1 登陆聊天室")
		fmt.Println("\t\t2 注册用户")
		fmt.Println("\t\t3 退出系统")
		fmt.Println("\t\t请选择(1-3)")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户的ID:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码:")
			fmt.Scanf("%s\n", &userPwd)
			// 完成登陆
			// 创建一个UserProcess实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
			//loop = false
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名称:")
			fmt.Scanf("%s\n", &userName)
			// 调用UserProcess实例完成注册请求
			up := &process.UserProcess{}
			up.Register(userId,userPwd,userName)
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
			//loop =false
		default :
			fmt.Println("你的输入有误，请重新输入:")
		}
	}
	//
	//// 根据用户的输入显示新的提示信息
	//if key == 1 {
	//	// 用户要登陆
	//
	//
	//	// 因为使用了分层结构，
	//
	//	// 先把登陆函数写到另外的文件
	//	// 这里需要重新调用
	//	// login(userId, userPwd)
	//	//if err != nil {
	//	//	fmt.Println("登陆失败!")
	//	//} else {
	//	//	fmt.Println("登陆成功!")
	//	//}
	//} else if key == 2 {
	//	fmt.Println("显示注册用户逻辑")
	//}
}
