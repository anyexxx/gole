package login

import (
	"me.qqtu.game/cache"
	"me.qqtu.game/message"
	"me.qqtu.game/model"
	"me.qqtu.game/notifier"
	"me.qqtu.game/server"
)

func HandleLogin(c server.Conn, msg message.Message) {
	var account string = "account1"
	var pwd string = "pwd1"

	//json test
	data := new(message.CGLoginData)
	msg.Decode(data)

	account = data.Account
	pwd = data.Pwd

	accModel := &model.AccountModel{Account: account, Password: pwd}
	//测试，若无账号，则直接注册
	if accModel.Login() == false {
		accModel.Register()
	}
	if accModel.UserId > 0 {
		u := new(server.User)
		u.Conn = c
		u.Id = accModel.UserId
		//创建用户数据缓存
		u.Data = &cache.UserData{UserId: u.Id}
		//创建用户关闭连接回调
		u.CloseHandler = func(u *server.User) {
			notifier.UserLogout(u.Id)
			server.QuitUser(u.Id)
		}
		c.SetUser(u)
		//管理
		server.AddUser(u)        //在线用户管理
		notifier.UserLogin(u.Id) //广播管理
	}
	//返回数据
	resData := new(message.GCLoginData)
	resData.Ret = 1
	c.SendMessage(message.NewJsonMessage(message.MSG_GC_LOGIN, resData))
}
