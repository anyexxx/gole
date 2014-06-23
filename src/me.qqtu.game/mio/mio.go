package mio

import (
	"labix.org/v2/mgo"
	"strings"
)

const (
	//DB name
	PLAYER_DB string = "playerdb"
	CONFIG_DB string = "config"
	//Collection name
	ROLE_C    string = "roles"
	ACCOUNT_C string = "accounts"

	GAMESERVER_C  string = "gameservers"
	REDISSERVER_C string = "redisservers"
)

var (
	mdb *MongoDB
)

func InitMongoDB(url string, user string, pwd string) {
	mdb = &MongoDB{addrs: strings.Split(url, ";"), user: user, pwd: pwd}
	mdb.Conn()
	//init db and collection
	GetRoleC().EnsureIndexKey("userid")
}

func Close() {
	if mdb != nil {
		mdb.Close()
	}
}

func GetPlayerDB() *mgo.Database {
	if mdb != nil {
		return mdb.session.DB(PLAYER_DB)
	}
	return nil
}

func GetConfigDB() *mgo.Database {
	if mdb != nil {
		return mdb.session.DB(CONFIG_DB)
	}
	return nil
}

func GetGameServerC() *mgo.Collection {
	return GetConfigDB().C(GAMESERVER_C)
}

func GetRedisServerC() *mgo.Collection {
	return GetConfigDB().C(REDISSERVER_C)
}

func GetRoleC() *mgo.Collection {
	return GetPlayerDB().C(ROLE_C)
}

func GetAccountC() *mgo.Collection {
	return GetPlayerDB().C(ACCOUNT_C)
}
