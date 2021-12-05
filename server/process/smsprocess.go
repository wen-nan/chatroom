package process

import (
	"encoding/json"
	"fmt"
	"net"
	"server/common/message"
	"server/utils"
)

type SmsProcess struct {
}

// 转发消息

func (this *SmsProcess)SendGroupMes(mes *message.Message) {
	// 遍历服务器端onlineUsers
	// 将消息转发

	// 取出mes的内容
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
	}

	for id, up := range userMgr.onlineUsers {
		// 过滤掉自己
		if id == smsMes.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}
}

func (this *SmsProcess)SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	// 创建一个tf实例，发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}