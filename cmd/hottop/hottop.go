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

const getTimelineLua = `
local key = KEYS[1]
local begin = ARGV[1]
local timeline = redis.call('GET', key)
if not timeline then
    redis.call('SET', key, begin)
    timeline = begin
end
return timeline
`

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
	timeline, err := redis.Int64(cli.Do("EVAL", getTimelineLua, 1, redisHotTopTimeLineKey(t), begin))
	if err != nil {
		fmt.Printf("get timeline error: %v\n", err)
		return begin
	}
	return timeline
}

func setHotTimeLine(newBeginTime int64, t string) {
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	if _, err := cli.Do("SET", redisHotTopTimeLineKey(t), newBeginTime); err != nil {
		fmt.Printf("set timeline error: %v\n", err)
	}
}

func parseData(t int, arr interface{}, year int) error {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("arr not slice")
	}
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	var hotItem []hotso.HotItem
	bytes, err := bson.MarshalJSON(arr)
	if err != nil {
		return err
	}
	if err := bson.UnmarshalJSON(bytes, &hotItem); err != nil {
		return err
	}

	hotType := strings.ToLower(hotso.HotSoType[t])
	detailKey := internal.GetHotDetailKey(hotType, year)
	topKey := internal.GetHotTopKey(hotType, year)

	for _, val := range hotItem {
		if val.State == "荐" {
			continue
		}
		var iRead int
		if val.Reading == "" {
			top, err := strconv.Atoi(val.Top)
			if err != nil {
				continue
			}
			iRead = (100 - top) * 100
		} else {
			reading, err := strconv.Atoi(val.Reading)
			if err != nil {
				continue
			}
			iRead = reading
		}
		titleMd5 := common.MD5(val.Title)
		b, err := json.Marshal(val)
		if err != nil {
			continue
		}
		mFields := map[string]string{titleMd5: string(b)}
		if _, err := cli.Do("HMSET", redis.Args{}.Add(detailKey).AddFlat(mFields)...); err == nil {
			score := fmt.Sprintf("%.3f", float32(iRead)/10000.0)
			if _, err := cli.Do("ZINCRBY", topKey, score, titleMd5); err != nil {
				fmt.Printf("zincrby error: %v\n", err)
			}
		}
	}
	return nil
}

func produceData(t int) {
	defer wg.Done()
	if val, ok := hotso.HotSoType[t]; ok {
		year := time.Now().Year()
		hotType := strings.ToLower(val)
		datas := internal.NewMongoDB().OnLoadData(t, getHotTimeLine(hotType), getEndTime(hotType))
		for _, v := range datas {
			if err := parseData(t, v.Data, year); err != nil {
				fmt.Printf("parse data error: %v\n", err)
			}
		}
		setHotTimeLine(getEndTime(hotType), hotType)
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
		for k := range hotso.HotSoType {
			go produceData(k)
		}
	}
	wg.Wait()
}
