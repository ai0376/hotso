package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mjrao/hotso/config"
	"github.com/mjrao/hotso/internal"
)

func nowTime() int64 {
	return time.Now().Unix()
}

func getBeginTime() int64 {
	return config.GetConfig().HotTop.BeginTime
}

func getDurationTimeSec() int64 {
	return config.GetConfig().HotTop.DurationTimeSec
}

func getEndTime() int64 {
	if getDurationTimeSec() <= 0 {
		return nowTime()
	}
	return getDurationTimeSec() + nowTime()
}

func redisHotTopTimeLineKey() string {
	return "hottoptimeline"
}

func getHotTimeLine() int64 {
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	begin := getBeginTime()
	if timeline, err := redis.Int64(cli.Do("GET", redisHotTopTimeLineKey())); err != nil {
		cli.Do("SET", redisHotTopTimeLineKey(), begin)
		timeline = begin
	} else {
		begin = timeline
	}
	return begin
}

func setHotTimeLine(newBeginTime int) {
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	if _, err := cli.Do("SET", redisHotTopTimeLineKey(), newBeginTime); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println(getHotTimeLine())
}
