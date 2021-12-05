package process

import (
	"client/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) { // 这个地方mes一定是SmsMes
	// 显示即可
	// 反序列化mes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	// 显示信息
	info := fmt.Sprintf("用户id:\t%d对大家说：\t%s\n", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
}
