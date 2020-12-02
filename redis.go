package heeglibs

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	Host   string
	Port   string
	PassWd string

	Db      int
	Timeout int

	rds *redis.Client
}

func NewRedis(host, port, passwd string, db, timeout int) *Redis {
	rds := &Redis{
		Host:    host,
		Port:    port,
		PassWd:  passwd,
		Db:      db,
		Timeout: timeout,
	}

	if 0 == timeout {
		rds.Timeout = 30 * 60
	}

	return rds
}

// 获取缓存在线用户的数据库
func (this *Redis) Open() {
	if nil == this.rds {
		this.rds = redis.NewClient(&redis.Options{
			Addr:     this.Host + ":" + this.Port,
			Password: this.PassWd, // no password set
			DB:       this.Db,     // use default DB
		})

		_, err := this.rds.Ping().Result()
		if nil != err {
			panic("Get rds Error:" + err.Error())
		}

		log.Println("Create rds success!")
	}

	return
}

// 存储在线用户缓存数据
//
// @param key 	缓存数据key
// @param value	缓存的数据
// @param timeout 设置超时时间
//
// @return err 	错误信息
func (this *Redis) SetValue(key string, value interface{}, timeout int) (err error) {
	if nil == this.rds {
		err = errors.New("Redis didn't open.")

		return
	}

	tout := this.Timeout
	if 0 < timeout {
		tout = timeout
	}

	err = this.rds.Set(key, value, time.Duration(tout)*time.Second).Err()
	if err != nil {
		log.Error(key, "Set Error", err.Error())
	}

	err = this.rds.Expire(key, time.Duration(tout)*time.Second).Err()
	if nil != err {
		log.Error(err)
	}

	return
}

// 获取在线用户缓存数据
//
// @param key 	要获取缓存数据的key
// @param isset 是否重新更新缓存时间
// @param timeout 重新设置更新的时间
//
// @return (value,err)
func (this *Redis) GetValue(key string, isset bool, timeout int) (value interface{}, err error) {
	if nil == this.rds {
		err = errors.New("Redis didn't open.")

		return
	}

	value, err = this.rds.Get(key).Result()
	if err == redis.Nil {
		log.Error("key does not exist")
	} else if err != nil {
		log.Error(err.Error())
	}

	tout := this.Timeout
	if 0 < timeout {
		tout = timeout
	}

	if isset {
		err = this.rds.Expire(key, time.Duration(tout)*time.Second).Err()
		if nil != err {
			log.Error(err)
		}
	}

	return
}

// 获取在线缓存的所有keys
//
// @return (keys[],err)
func (this *Redis) GetKeys() (keys []string, err error) {
	if nil == this.rds {
		err = errors.New("Redis didn't open.")

		return
	}

	ks := this.rds.Keys("*")
	keys, err = ks.Result()

	return
}

// 删除在线的用户缓存
//
// @param key 	要删除缓存的key
func (this *Redis) DelValue(key string) {
	if nil == this.rds {
		return
	}

	this.rds.Del(key).Result()
}

// 清除用户缓存中的所有数据
func (this *Redis) FlushAll() {
	keys, err := this.GetKeys()
	if nil == err {
		for _, key := range keys {
			this.DelValue(key)
		}
	}
}
