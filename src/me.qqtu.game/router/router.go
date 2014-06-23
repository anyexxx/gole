package router

import (
	"me.qqtu.game/logger"
	"me.qqtu.game/mio"
	"time"
)

const (
	ID_BASE int64 = 10000000000 //用户id为此基数乘以服务器id值，再加上当前服务器对应自增数，用以区分不同redis
)

var (
	CurrSID      int                  //当前服务器id
	CurrHost     string               //当前服务器ip
	CurrPort     string               //当前服务端口
	gameServers  map[int]*DestServer  //游戏服务器列表
	redisServers map[int]*RedisServer //redis列表

	lastloadTime       int64 //上次拉取配置的时间戳
	loadconfigInterval int64 //拉取配置的最小时间间隔
)

func InitRouter(currHost string, currPort string) {
	gameServers = make(map[int]*DestServer)
	redisServers = make(map[int]*RedisServer)
	//for test
	CurrHost = currHost
	CurrPort = currPort
	loadconfigInterval = 60 //最小间隔默认设置为60秒

	//测试设置
	gameServers[1] = &DestServer{Id: 1, Host: "127.0.0.1", Port: "2000"}
	redisServers[1] = &RedisServer{DestServer: DestServer{Id: 1, Host: "127.0.0.1", Port: "6379"}, RouterUidPerfix: []int{1}}

	//加载配置
	LoadConfig()
	CurrSID = GetGameServerByRemote(CurrHost, CurrPort).Id

	//
	logger.GetLogger().Info("Init router done!", logger.Extras{"currHost": currHost, "currPort": currPort, "gs": gameServers, "rs": redisServers})
}

type DestServer struct {
	Id   int
	Host string
	Port string
}

type RedisServer struct {
	DestServer
	RouterUidPerfix []int //对于给定的一个uid，检测是否该连接到此redis
}

//检查是否匹配redis路由规则
func (s *RedisServer) MatchRouter(userId int64) bool {
	routerId := int(userId / ID_BASE)
	if routerId <= 0 {
		logger.GetLogger().Error("router.MatchRouter error: invaild userid", logger.Extras{"uid": userId})
		return false
	}
	for _, v := range s.RouterUidPerfix {
		if v == routerId {
			return true
		}
	}
	return false
}

func LoadConfig() {
	//每次拉取有一定时间间隔
	nowTime := time.Now().Unix()
	if nowTime-lastloadTime < loadconfigInterval {
		return
	}
	lastloadTime = nowTime

	//mongodb存入的是bson格式数据
	var gs []DestServer
	var rs []RedisServer
	mio.GetGameServerC().Find(nil).All(&gs)
	mio.GetRedisServerC().Find(nil).All(&rs)

	gameServers = make(map[int]*DestServer)
	for _, v := range gs {
		gameServers[v.Id] = &v
	}

	redisServers = make(map[int]*RedisServer)
	for _, v := range rs {
		redisServers[v.Id] = &v
	}
}

func ReloadConfig() {
	lastloadTime = 0
	LoadConfig()
}

func GetCurrGameServer() *DestServer {
	return GetGameServerByRemote(CurrHost, CurrPort)
}

func GetGameServerByRemote(host string, port string) *DestServer {
	for i := 0; i < 2; i++ {
		for _, v := range gameServers {
			if v != nil && v.Port == port && v.Host == host {
				return v
			}
		}
		if i == 0 {
			LoadConfig()
		}
	}
	return nil
}

func GetGameServerById(id int) *DestServer {
	for i := 0; i < 2; i++ {
		for _, v := range gameServers {
			if v != nil && v.Id == id {
				return v
			}
		}
		if i == 0 {
			LoadConfig()
		}
	}
	return nil
}

func GetReidsServerByUserId(userId int64) *RedisServer {
	for i := 0; i < 2; i++ {
		for _, v := range redisServers {
			if v != nil && v.MatchRouter(userId) {
				return v
			}
		}
		if i == 0 {
			LoadConfig()
		}
	}
	return nil
}

func GetRedisServerById(id int) *RedisServer {
	for i := 0; i < 2; i++ {
		for _, v := range redisServers {
			if v != nil && v.Id == id {
				return v
			}
		}
		if i == 0 {
			LoadConfig()
		}
	}
	return nil
}

func GetReidsServerIdByUserId(userId int64) int {
	ret := GetReidsServerByUserId(userId)
	if ret == nil {
		return -1
	}
	return ret.Id
}
