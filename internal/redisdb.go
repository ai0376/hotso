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
	hotdataKey = "hotword:%s:%d"
)

//GetHotWordKey ...
func GetHotWordKey(hot_type string, year int) string {
	return fmt.Sprintf(hotdataKey, hot_type, year)
}
