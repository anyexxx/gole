package handler

import (
	"me.qqtu.game/handler/login"
	"me.qqtu.game/handler/test"
	"me.qqtu.game/message"
	"me.qqtu.game/server"
)

func InitHandlers() {
	server.RegisterHandler(message.MSG_TEST, test.HandleTest)
	server.RegisterHandler(message.MSG_CG_LOGIN, login.HandleLogin)
}
