package notifier

import (
	"fmt"
	"me.qqtu.game/message"
	"me.qqtu.game/router"
	"me.qqtu.game/server"
	"net"
)

var (
	client *server.SocketClient
	host   string
	port   string
)

func InitNotifierServer(host string, port string) {
	host = host
	port = port
	client = server.CreateSocketClient(host, port, message.DATAFORMAT_JSON)
	client.Run()
}

func Close() {
	if client != nil {
		client.Conn.Close()
	}
}

func sendMessage(msg message.Message) {
	if client != nil && client.Conn != nil {
		client.Conn.SendMessage(msg)
	}
}

func Broadcast(msg message.Message) {
	sendMessage(msg)
}

func ServerLogin() {
	data := &message.NotiServerLogin{SId: router.CurrSID}
	msg := message.NewJsonMessage(message.MSG_NOTI_SERVER_LOGIN, data)
	sendMessage(msg)
}

func UserLogin(uid int64) {
	data := &message.NotiUserLogin{SId: router.CurrSID, UId: uid}
	msg := message.NewJsonMessage(message.MSG_NOTI_USER_LOGIN, data)
	sendMessage(msg)
}

func UserLogout(uid int64) {
	data := &message.NotiUserLogout{SId: router.CurrSID, UId: uid}
	msg := message.NewJsonMessage(message.MSG_NOTI_USER_LOGOUT, data)
	sendMessage(msg)
}

//添加广播消息处理办法
func handleNofify(c net.Conn, msg message.Message) {
	switch msg.MessageId() {
	case message.MSG_NOTI_BROADCAST:
		fmt.Println("broadcast:", msg.MessageId())
	}
}
