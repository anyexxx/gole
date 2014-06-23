package server

import (
	"me.qqtu.game/message"
)

var (
	msgMap map[uint16]func(Conn, message.Message) //保存消息id和相应处理handler的对应关系
)

//初始化
func init() {
	if msgMap == nil {
		msgMap = make(map[uint16]func(Conn, message.Message))
	}
}

//注册消息处理
func RegisterHandler(msgId uint16, handler func(Conn, message.Message)) {
	msgMap[msgId] = handler
}

//获取消息处理
func GetMsgHandler(msgId uint16) func(Conn, message.Message) {
	return msgMap[msgId]
}

//检查是否有此消息的处理方法
func HasMsgHandler(msgId uint16) bool {
	_, prs := msgMap[msgId]
	return prs
}

//取消消息处理
func UnRegisterHandler(msgId uint16) {
	delete(msgMap, msgId)
}
