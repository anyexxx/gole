package main

import (
	"code.google.com/p/goconf/conf"
	"me.qqtu.game/logger"
	"me.qqtu.game/message"
	"me.qqtu.game/server"
	"os"
	"runtime"
)

var (
	usMap  map[int]map[int64]bool //服务器与用户对照表
	scMap  map[int]server.Conn    //服务器id与连接对照表
	config = ` [server]
				host=127.0.0.1
				port=2001
				`
)

func main() {
	//设置cpu数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg, _ := conf.ReadConfigBytes([]byte(config))
	host, _ := cfg.GetString("server", "host")
	port, _ := cfg.GetString("server", "port")

	usMap = make(map[int]map[int64]bool)
	scMap = make(map[int]server.Conn)

	//注册服务响应
	server.RegisterHandler(message.MSG_NOTI_BROADCAST, HandleBroadcast)
	server.RegisterHandler(message.MSG_NOTI_SERVER_LOGIN, HandleGameServerConn)
	server.RegisterHandler(message.MSG_NOTI_SERVER_LOGOUT, HandleGameServerDisConn)
	server.RegisterHandler(message.MSG_NOTI_USER_LOGIN, HandleUserLogin)
	server.RegisterHandler(message.MSG_NOTI_USER_LOGOUT, HandleUserLogout)

	logger.GetLogger().Info("Notifier Server Start!", logger.Extras{"pid": os.Getpid()})

	sc := server.CreateSocketServer(host, port, message.DATAFORMAT_JSON)
	sc.Run()
}

//游戏服务器连接
func HandleGameServerConn(c server.Conn, msg message.Message) {
	var sid int = 1
	scMap[sid] = c
	logger.GetLogger().Info("Notifier: server connect!", logger.Extras{"sid": sid})
}

//游戏服务器断开
func HandleGameServerDisConn(c server.Conn, msg message.Message) {
	var sid int = 1
	if _, prs := scMap[sid]; prs {
		delete(scMap, sid)
	}
	logger.GetLogger().Info("Notifier: server disconnect!", logger.Extras{"sid": sid})
}

//对指定用户群发送消息
func HandleSendMsgToUser(c server.Conn, msg message.Message) {
	var users []int64
	for _, uid := range users {
		for k, v := range usMap {
			var sid int = k
			if _, prs := v[uid]; prs {
				if conn := scMap[sid]; conn != nil {
					conn.SendMessage(msg)
					logger.GetLogger().Info("Notifier: send msg to user!", logger.Extras{"msgId": msg.MessageId(), "uid": uid})
				}
			}
		}
	}
}

//广播消息
func HandleBroadcast(c server.Conn, msg message.Message) {
	for _, conn := range scMap {
		if conn != nil {
			conn.SendMessage(msg)
			logger.GetLogger().Info("Notifier: broadcast message.", logger.Extras{"msgId": msg.MessageId()})
		}
	}
}

//用户登录注册
func HandleUserLogin(c server.Conn, msg message.Message) {
	var userId int64 = 100
	var serverId int = 1

	if _, prs := usMap[serverId]; prs == false {
		usMap[serverId] = make(map[int64]bool)
	}
	usMap[serverId][userId] = true
	logger.GetLogger().Info("Notifier: user login.", logger.Extras{"sid": serverId, "uid": userId})
}

//用户退出
func HandleUserLogout(c server.Conn, msg message.Message) {
	var userId int64 = 100

	for _, v := range usMap {
		if _, prs := v[userId]; prs {
			delete(v, userId)
		}
	}
	logger.GetLogger().Info("Notifier: user logout.", logger.Extras{"uid": userId})
}
