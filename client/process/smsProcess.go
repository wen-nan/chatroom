package process

import (
	"client/common/message"
	"client/utils"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

// 发送群聊消息

func (this *SmsProcess)SendGroupMes(content string) (err error) {
	// 创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	// 创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content // 内容
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	// 序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err.Error())
		return
	}
	mes.Data = string(data)

	// 对mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err.Error())
		return
	}

	// 将mes发送个服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	// 发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes err=", err.Error())
	}
	return
}