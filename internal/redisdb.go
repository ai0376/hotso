package internal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mjrao/hotso/config"
)

/*
useage:

RedisCliPool()

*/
var cliPool *redis.Pool

//NewRedisCliPool ...
func NewRedisCliPool(maxIdle, maxActive, idleTimeOut int, host string, port int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeOut) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+strconv.Itoa(port))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

//RedisCliPool ...
func RedisCliPool() *redis.Pool {
	if cliPool != nil {
		return cliPool
	}
	cliPool = NewRedisCliPool(10, 100, 20, config.GetConfig().Redis.Host, config.GetConfig().Redis.Port)
	return cliPool
}

const (
	hotDataKey      = "hotword:%d:%s"   //hotword:2019:weibo
	hotTopDataKey   = "hottop:%d:%s"    //hottop:2019:weibo
	hotTopDetailKey = "hotdetail:%d:%s" //hotdetail:2019:weibo
)

//GetHotWordKey ...
func GetHotWordKey(hotType string, year int) string {
	return fmt.Sprintf(hotDataKey, year, hotType)
}

//GetHotTopKey ...
func GetHotTopKey(hotType string, year int) string {
	return fmt.Sprintf(hotTopDataKey, year, hotType)
}

//GetHotDetailKey ...
func GetHotDetailKey(hotType string, year int) string {
	return fmt.Sprintf(hotTopDetailKey, year, hotType)
}
