package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/mjrao/hotso/config"
	"github.com/mjrao/hotso/internal"
	"github.com/mjrao/hotso/internal/metadata/hotso"
)

//timeOfDay day is "2006-01-02"
func timeOfDay(day string) (int64, int64) {
	timeStrBegin := day + " 00:00:01"
	timeStrEnd := day + " 23:59:59"

	t1, err1 := time.Parse("2006-01-02 15:04:05", timeStrBegin)
	t2, err2 := time.Parse("2006-01-02 15:04:05", timeStrEnd)
	if err1 != nil || err2 != nil {
		return -1, -1
	}
	return t1.Unix() - 8*3600, t2.Unix() - 8*3600
}

//ResponIndentJSON ...
func ResponIndentJSON(c *gin.Context, code int, obj interface{}) {
	w := c.Writer
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
	}
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u003c"), []byte("<"), -1)
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u003e"), []byte(">"), -1)
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u0026"), []byte("\u0026"), -1)
	c.Status(code)
	w.Write(jsonBytes)
}

func convertHotData(data *hotso.HotData, num int) (*hotso.HotData, error) {
	if data == nil {
		return nil, fmt.Errorf("no data found")
	}
	dataSlice, ok := data.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}
	var resultData []map[string]interface{}
	for i, v := range dataSlice {
		if num != 0 && i >= num {
			break
		}
		if m, ok := v.(map[string]interface{}); ok {
			resultData = append(resultData, m)
		}
	}
	return &hotso.HotData{Type: data.Type, Name: data.Name, InTime: data.InTime, Data: resultData}, nil
}

//GetDataByType ...
func GetDataByType(dataType int, num int) (*hotso.HotData, error) {
	data := internal.NewMongoDB().OnFindOneDataByType(dataType)
	return convertHotData(data, num)
}

//GetHotWordData ...
func GetHotWordData(c *gin.Context) {
	var hottype string
	switch c.Param("hottype") {
	case "weibo":
		hottype = "weibo"
	case "baidu":
		hottype = "baidu"
	default:
		c.JSON(http.StatusOK, gin.H{"code": -1, "message": "no data"})
		return
	}
	num, _ := strconv.Atoi(c.Param("num"))
	year, _ := strconv.Atoi(c.Param("year"))
	key := internal.GetHotWordKey(hottype, year)
	if num <= 0 || num > 100 {
		num = 100
	}
	var result []map[string]interface{}
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	reply, err := redis.Values(cli.Do("ZREVRANGE", key, 0, num-1, "WITHSCORES"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	for i := 0; i < len(reply); i += 2 {
		result = append(result, map[string]interface{}{
			"rank":  i/2 + 1,
			"word":  string(reply[i].([]byte)),
			"score": string(reply[i+1].([]byte)),
		})
	}
	ResponIndentJSON(c, http.StatusOK, result)
}

//GetHotTopData 年度总榜
func GetHotTopData(c *gin.Context) {
	var hottype string
	switch c.Param("hottype") {
	case "weibo":
		hottype = "weibo"
	case "baidu":
		hottype = "baidu"
	case "zhihu":
		hottype = "zhihu"
	case "shuimu":
		hottype = "shuimu"
	case "tianya":
		hottype = "tianya"
	case "v2ex":
		hottype = "v2ex"
	default:
		c.JSON(http.StatusOK, gin.H{"code": -1, "message": "no data"})
		return
	}
	num, _ := strconv.Atoi(c.Param("num"))
	year, _ := strconv.Atoi(c.Param("year"))
	key := internal.GetHotTopKey(hottype, year)
	if num <= 0 || num > 100 {
		num = 100
	}
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	reply, err := redis.Values(cli.Do("ZREVRANGE", key, 0, num-1, "WITHSCORES"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	var array []string
	for i := 0; i < len(reply); i += 2 {
		array = append(array, string(reply[i].([]byte)))
	}
	var result []hotso.HotItem
	if len(array) > 0 {
		reply2, err := redis.Values(cli.Do("HMGET", redis.Args{}.Add(internal.GetHotDetailKey(hottype, year)).AddFlat(array)...))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
			return
		}
		for i, v := range reply2 {
			if v == nil {
				continue
			}
			var item hotso.HotItem
			if err := json.Unmarshal(v.([]byte), &item); err != nil {
				continue
			}
			item.Top = strconv.Itoa(i + 1)
			result = append(result, item)
		}
	}
	ResponIndentJSON(c, http.StatusOK, result)
}

//GetHotType ...
func GetHotType(c *gin.Context) {
	num, _ := strconv.Atoi(c.Param("num"))
	var data *hotso.HotData
	var err error
	switch c.Param("hottype") {
	case "weibo":
		data, err = GetDataByType(hotso.WEIBO, num)
	case "baidu":
		data, err = GetDataByType(hotso.BAIDU, num)
	case "zhihu":
		data, err = GetDataByType(hotso.ZHIHU, num)
	case "shuimu":
		data, err = GetDataByType(hotso.SHUIMU, num)
	case "tianya":
		data, err = GetDataByType(hotso.TIANYA, num)
	case "v2ex":
		data, err = GetDataByType(hotso.V2EX, num)
	default:
		c.JSON(http.StatusOK, gin.H{"code": -1, "message": "no data"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	ResponIndentJSON(c, http.StatusOK, data)
}

//QueryDataOfDay ...
func QueryDataOfDay(dataType int, num int, day string) (*hotso.HotData, error) {
	start, end := timeOfDay(day)
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("invalid day format")
	}
	if _, ok := hotso.HotSoType[dataType]; !ok {
		return nil, fmt.Errorf("invalid hot type")
	}
	data := internal.NewMongoDB().OnQueryData(dataType, start, end)
	return convertHotData(data, num)
}

//QueryHotSoOfDay ...
func QueryHotSoOfDay(c *gin.Context) {
	day := c.Param("day")
	num, _ := strconv.Atoi(c.Param("num"))
	var data *hotso.HotData
	var err error
	switch c.Param("hottype") {
	case "weibo":
		data, err = QueryDataOfDay(hotso.WEIBO, num, day)
	case "baidu":
		data, err = QueryDataOfDay(hotso.BAIDU, num, day)
	case "zhihu":
		data, err = QueryDataOfDay(hotso.ZHIHU, num, day)
	case "shuimu":
		data, err = QueryDataOfDay(hotso.SHUIMU, num, day)
	case "tianya":
		data, err = QueryDataOfDay(hotso.TIANYA, num, day)
	case "v2ex":
		data, err = QueryDataOfDay(hotso.V2EX, num, day)
	default:
		c.JSON(http.StatusOK, gin.H{"code": -1, "message": "no data"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": err.Error()})
		return
	}
	ResponIndentJSON(c, http.StatusOK, data)
}

func main() {
	serviceCfg := config.GetConfig().Service
	router := gin.Default()
	v1 := router.Group("hotso/v1")
	{
		v1.GET("/hotso/:hottype/:num", GetHotType)
		v1.GET("/hotword/:hottype/:year/:num", GetHotWordData)
		v1.GET("/hottop/:hottype/:year/:num", GetHotTopData)
		v1.GET("/query/:hottype/:day/:num", QueryHotSoOfDay)
	}
	addr := fmt.Sprintf("%s:%d", serviceCfg.IP, serviceCfg.Port)
	router.Run(addr)
}
