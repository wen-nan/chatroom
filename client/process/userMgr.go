package process

import (
	"client/common/message"
	"client/model"
	"fmt"
)

// 客户端要维护的map
var onlineUsers map[int]*model.User= make(map[int]*model.User, 10)
// 因为在客户端很多地方会用到curUser,用一个全局变量

var CurUser model.CurUser  // 用户登陆成功后，完成初始化

// 客户端显示当前在线用户
func outputOnlineUser() {
	// 遍历onlineUsers
	fmt.Println("当前在线用户列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户id:", id)
	}
}

// 编写一个方法，处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	// 优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok { // 原来没有
		user = &model.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	outputOnlineUser()
}