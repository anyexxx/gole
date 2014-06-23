package rio

import (
	"errors"
	"me.qqtu.game/logger"
	"me.qqtu.game/router"
	"time"
)

//根据router中的配置，为每个user获取相应的redis连接

var (
	dbs map[int]*DBRedis
)

func init() {
	dbs = make(map[int]*DBRedis)
}

func GetRedisByUserId(userId int64) *DBRedis {
	destServer := router.GetReidsServerByUserId(userId)
	if destServer != nil {
		//获取已缓存的redis连接
		ret, prs := dbs[destServer.Id]
		if prs {
			return ret
		}
		//创建新的redis连接
		dbRedis := &DBRedis{
			id:       destServer.Id,
			host:     destServer.Host,
			port:     destServer.Port,
			poolsize: 5,
			timeout:  240 * time.Second}
		if err := dbRedis.Conn(); err != nil {
			logger.GetLogger().Error("redis connect error!", logger.Extras{"error": err})
			return nil
		}
		dbs[destServer.Id] = dbRedis
		return dbRedis
	}
	logger.GetLogger().Error("no redis match this uid", logger.Extras{"uid": userId})
	return nil
}

func GetReidsByServerId(sId int) *DBRedis {
	//获取一个用户id范围，用以确定连接的redis
	return GetRedisByUserId(getRealIncrId(1, sId))
}

func Close(id int) error {
	dbRedis, ok := dbs[id]
	if !ok {
		err := errors.New("close nil id")
		logger.GetLogger().Error(err.Error(), logger.Extras{"id": id, "err": err})
		return err
	}
	if dbRedis.pool != nil {
		defer dbRedis.pool.Close() //确定关闭
	}
	delete(dbs, id)
	return nil
}

func CloseAll() {
	for k, _ := range dbs {
		Close(k)
	}
}
