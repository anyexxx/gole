package server

import (
	"me.qqtu.game/cache"
	"me.qqtu.game/logger"
)

var (
	users = make(map[int64]*User)
)

//
//当前服务器在线用户管理
//

func AddUser(u *User) bool {
	if u == nil {
		logger.GetLogger().Warning("server.AddUser, nil param", nil)
		return false
	}
	if _, prs := users[u.Id]; prs {
		logger.GetLogger().Warning("server.AddUser, user already added", logger.Extras{"uid": u.Id})
		return false
	}
	users[u.Id] = u
	return true
}

func GetUser(id int64) *User {
	if u, prs := users[id]; prs {
		return u
	}
	return nil
}

func NumberOfUser() int {
	return len(users)
}

func QuitUser(id int64) {
	if u, prs := users[id]; prs {
		u.Conn.Close()
		delete(users, id)
		logger.GetLogger().Info("User quit.", logger.Extras{"uid": id})
		return
	}
}

type User struct {
	Id           int64           //当前用户id
	Conn         Conn            //当前连接
	Data         *cache.UserData //用户保存的数据
	CloseHandler func(*User)     //连接关闭时回调函数
}
