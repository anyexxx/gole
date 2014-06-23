package mio

import (
	"labix.org/v2/mgo"
	"me.qqtu.game/logger"
	"strings"
)

type MongoDB struct {
	addrs   []string
	user    string
	pwd     string
	session *mgo.Session
}

func (m *MongoDB) Conn() {
	session, err := mgo.Dial(strings.Join(m.addrs, ",")) //连接mongodb集群
	if err != nil {
		panic(err)
	}
	//登录验证
	if m.user != "" && m.pwd != "" {
		session.Login(&mgo.Credential{Username: m.user, Password: m.pwd})
	}
	m.session = session
}

func (m *MongoDB) Close() {
	defer logger.GetLogger().HandlePanic()
	if m.session != nil {
		m.session.Close()
	}
}
