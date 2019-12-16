package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mjrao/hotso/common"
	"github.com/mjrao/hotso/config"
	"github.com/mjrao/hotso/internal"
	"github.com/mjrao/hotso/internal/metadata/hotso"
	"gopkg.in/mgo.v2/bson"
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

func getEndTime(t string) int64 {
	if getDurationTimeSec() <= 0 {
		return nowTime()
	}
	return getDurationTimeSec() + getHotTimeLine(t)
}

func redisHotTopTimeLineKey(t string) string {
	return fmt.Sprintf("hottoptimeline:%s", t)
}

func getHotTimeLine(t string) int64 {
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	begin := getBeginTime()
	if timeline, err := redis.Int64(cli.Do("GET", redisHotTopTimeLineKey(t))); err != nil {
		cli.Do("SET", redisHotTopTimeLineKey(t), begin)
		timeline = begin
	} else {
		begin = timeline
	}
	return begin
}

func setHotTimeLine(newBeginTime int64, t string) {
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	if _, err := cli.Do("SET", redisHotTopTimeLineKey(t), newBeginTime); err != nil {
		panic(err)
	}
}

func parseData(t int, arr interface{}) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("arr not slice")
	}
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	var hotItem []hotso.HotItem
	if bytes, err := bson.MarshalJSON(arr); err != nil {
		panic(err.Error())
	} else {
		bson.UnmarshalJSON(bytes, &hotItem)
	}
	var iRead = 0
	for _, val := range hotItem {
		if val.Reading == "" {
			top, _ := strconv.Atoi(val.Top)
			iRead = (10000 - top) * 100 //模拟出一个值
		} else {
			top, _ := strconv.Atoi(val.Reading)
			iRead = top
		}
		titleMd5 := common.MD5(val.Title)
		if b, err := json.Marshal(val); err == nil {
			mFields := map[string]string{
				titleMd5: string(b),
			}
			if _, err := cli.Do("HMSET", redis.Args{}.Add(internal.GetHotDetailKey(strings.ToLower(hotso.HotSoType[t]), time.Now().Year())).AddFlat(mFields)...); err == nil {
				score := fmt.Sprintf("%.3f", float32(iRead)/10000.0)
				cli.Do("ZINCRBY", internal.GetHotTopKey(strings.ToLower(hotso.HotSoType[t]), time.Now().Year()), score, titleMd5)
			}
		}
		return
	}
}

func produceData(t int) {
	defer wg.Done()
	if val, ok := hotso.HotSoType[t]; ok {
		datas := internal.NewMongoDB().OnLoadData(t, getHotTimeLine(strings.ToLower(val)), getEndTime(val))
		for _, v := range datas {
			parseData(t, v.Data)
			return
		}
		setHotTimeLine(getEndTime(val), strings.ToLower(val))
	}
}

var wg *sync.WaitGroup

func main() {
	wg = &sync.WaitGroup{}
	if len(os.Args) > 1 {
		for _, v := range os.Args[1:] {
			if n, err := strconv.Atoi(v); err != nil {
				fmt.Println("strconv Atoi error")
			} else {
				wg.Add(1)
				go produceData(n)
			}
		}
	} else {
		wg.Add(len(hotso.HotSoType))
		for k, _ := range hotso.HotSoType {
			go produceData(k)
		}
	}
	wg.Wait()
}
