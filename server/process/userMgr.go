package process

import "fmt"

// 因为userMgr实例在服务器端有且只有一个，并且在很多地方能用到，
// 所以定义为一个全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 完成初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 完成对onlineUsers添加

func (this *UserMgr)AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

// 删除

func (this *UserMgr)DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

// 返回当前所有在线用户

func (this *UserMgr)GetAllOnlineUsers() map[int]*UserProcess {
	return this.onlineUsers
}

// 根据id返回对应值

func (this *UserMgr)GetOnlineById(userId int) (up *UserProcess, err error) {
	// 如何从map中取出某个值
	up, ok := this.onlineUsers[userId]
	if !ok {
		// 说明要查找的用户当前不在线
		err = fmt.Errorf("用户%d不存在", userId)
		return
	}
	return
}
