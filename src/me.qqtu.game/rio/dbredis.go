package rio

import (
	"code.google.com/p/goprotobuf/proto"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"me.qqtu.game/logger"
	"me.qqtu.game/router"
	"time"
)

const (
	USER_ID_INCR_KEY string = "userIdIncr"
)

type DBRedis struct {
	id       int
	host     string
	port     string
	poolsize int
	timeout  time.Duration
	pool     *redis.Pool
}

func (r *DBRedis) Conn() error {
	host := r.host
	port := r.port
	poolsize := r.poolsize
	timeout := r.timeout

	r.pool = &redis.Pool{
		MaxActive:   poolsize,
		IdleTimeout: timeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+port)
			if err != nil {
				logger.GetLogger().Error(err.Error(), logger.Extras{"host": host, "port": port})
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func getRealIncrIdKey(key string, fixid int) string {
	return fmt.Sprint(key, ":", fixid)
}

func getRealIncrId(id int64, fixid int) int64 {
	fmt.Println(router.ID_BASE * int64(fixid))
	return router.ID_BASE*int64(fixid) + id
}

//生成用户数据存放的真正key
func GenUserKey(userId int64, key string) string {
	return fmt.Sprint(userId, ":", key)
}

func (r *DBRedis) GetActiveCount() int {
	return r.pool.ActiveCount()
}

func (r *DBRedis) GetConn() redis.Conn {
	return r.pool.Get()
}

func (r *DBRedis) Close() error {
	defer logger.GetLogger().HandlePanic()

	if r.pool != nil {
		return r.pool.Close()
	}
	return errors.New("redis pool is not initialized")
}

func (r *DBRedis) GetNextUserId() (int64, error) {
	return r.IncrIdVal(USER_ID_INCR_KEY)
}

func (r *DBRedis) GetUserIdIncr() (int64, error) {
	return r.GetIncrIdVal(USER_ID_INCR_KEY)
}

func (r *DBRedis) GetIncrIdVal(key string) (int64, error) {
	return redis.Int64(r.Get(getRealIncrIdKey(key, r.id)))
}

func (r *DBRedis) IncrIdVal(key string) (int64, error) {
	ret, err := r.Incr(getRealIncrIdKey(key, r.id))
	ret = getRealIncrId(ret, r.id)
	//
	logger.GetLogger().Error(err.Error(), nil)
	return ret, err
}

func (r *DBRedis) Pack(key string, val proto.Message) error {
	var err error
	if val != nil && key != "" {
		var bytes []byte
		bytes, err = proto.Marshal(val)
		if err != nil {
			logger.GetLogger().Error(err.Error(), nil)
			return err
		}
		_, err = r.GetConn().Do("SET", bytes)
		if err != nil {
			logger.GetLogger().Error(err.Error(), nil)
			return err
		}
	}
	return err
}

func (r *DBRedis) UnPack(key string, pb proto.Message) error {
	ret, err := r.Get(key)
	if err != nil {
		logger.GetLogger().Error(err.Error(), nil)
		return err
	}
	bytes := ret.([]byte)
	return proto.Unmarshal(bytes, pb)
}

func (r *DBRedis) Get(key string) (interface{}, error) {
	return r.GetConn().Do("GET", key)
}

func (r *DBRedis) Incr(key string) (int64, error) {
	return redis.Int64(r.GetConn().Do("INCR", key))
}

func (r *DBRedis) IncrBy(key string, inc int64) (int64, error) {
	return redis.Int64(r.GetConn().Do("INCRBY", key, inc))
}

func (r *DBRedis) Decr(key string) (int64, error) {
	return redis.Int64(r.GetConn().Do("DECR", key))
}

func (r *DBRedis) DecrBy(key string, dec int64) (int64, error) {
	return redis.Int64(r.GetConn().Do("DECRBY", key, dec))
}

func (r *DBRedis) HashSet(key string, sKey string, val interface{}) (bool, error) {
	ret, err := redis.Int(r.GetConn().Do("HSET", sKey, val))
	if err != nil {
		logger.GetLogger().Error(err.Error(), nil)
		return false, err
	} else {
		if ret == 1 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}
