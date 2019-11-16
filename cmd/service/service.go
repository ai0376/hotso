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
	"gopkg.in/mgo.v2/bson"
)

//ResponIndentJSON ...
func ResponIndentJSON(c *gin.Context, code int, obj interface{}) error {
	c.Status(code)
	w := c.Writer
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
	}
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u003c"), []byte("<"), -1)
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u003e"), []byte(">"), -1)
	jsonBytes = bytes.Replace(jsonBytes, []byte("\\u0026"), []byte("&"), -1)
	_, err = w.Write(jsonBytes)
	return err
}

//GetWeiBoData ...
func GetWeiBoData(num int) *hotso.HotData {
	data := internal.NewMongoDB().OnWeiBoFindOne()
	var hotdata hotso.HotData
	if bytes, err := bson.MarshalJSON(data); err != nil {
		panic(err.Error())
	} else {
		bson.UnmarshalJSON(bytes, &hotdata)
	}
	var resultData []map[string]interface{}
	index := 0
	for _, v := range hotdata.Data.([]interface{}) {
		index++
		if num != 0 && index > num {
			break
		}
		resultData = append(resultData, v.(map[string]interface{}))
	}
	return &hotso.HotData{Type: hotdata.Type, Name: hotdata.Name, InTime: hotdata.InTime, Data: resultData}
}

//GetBaiDuData ...
func GetBaiDuData(num int) *hotso.HotData {
	data := internal.NewMongoDB().OnBaiDuFindOne()
	var hotdata hotso.HotData
	if bytes, err := bson.MarshalJSON(data); err != nil {
		panic(err.Error())
	} else {
		bson.UnmarshalJSON(bytes, &hotdata)
	}
	var resultData []map[string]interface{}
	index := 0
	for _, v := range hotdata.Data.([]interface{}) {
		index++
		if num != 0 && index > num {
			break
		}
		resultData = append(resultData, v.(map[string]interface{}))
	}
	return &hotso.HotData{Type: hotdata.Type, Name: hotdata.Name, InTime: hotdata.InTime, Data: resultData}
}

//GetZhiHuData ...
func GetZhiHuData(num int) *hotso.HotData {
	data := internal.NewMongoDB().OnZhiHuFindOne()
	var hotdata hotso.HotData
	if bytes, err := bson.MarshalJSON(data); err != nil {
		panic(err.Error())
	} else {
		bson.UnmarshalJSON(bytes, &hotdata)
	}
	var resultData []map[string]interface{}
	index := 0
	for _, v := range hotdata.Data.([]interface{}) {
		index++
		if num != 0 && index > num {
			break
		}
		resultData = append(resultData, v.(map[string]interface{}))
	}
	return &hotso.HotData{Type: hotdata.Type, Name: hotdata.Name, InTime: hotdata.InTime, Data: resultData}
}

//GetShuiMuData ...
func GetShuiMuData(num int) *hotso.HotData {
	data := internal.NewMongoDB().OnShuiMuFindOne()
	var hotdata hotso.HotData
	if bytes, err := bson.MarshalJSON(data); err != nil {
		panic(err.Error())
	} else {
		bson.UnmarshalJSON(bytes, &hotdata)
	}
	var resultData []map[string]interface{}
	index := 0
	for _, v := range hotdata.Data.([]interface{}) {
		index++
		if num != 0 && index > num {
			break
		}
		resultData = append(resultData, v.(map[string]interface{}))
	}
	return &hotso.HotData{Type: hotdata.Type, Name: hotdata.Name, InTime: hotdata.InTime, Data: resultData}
}

//GetTianYaData ...
func GetTianYaData(num int) *hotso.HotData {
	data := internal.NewMongoDB().OnTianYaFindOne()
	var hotdata hotso.HotData
	if bytes, err := bson.MarshalJSON(data); err != nil {
		panic(err.Error())
	} else {
		bson.UnmarshalJSON(bytes, &hotdata)
	}
	var resultData []map[string]interface{}
	index := 0
	for _, v := range hotdata.Data.([]interface{}) {
		index++
		if num != 0 && index > num {
			break
		}
		resultData = append(resultData, v.(map[string]interface{}))
	}
	return &hotso.HotData{Type: hotdata.Type, Name: hotdata.Name, InTime: hotdata.InTime, Data: resultData}
}

//GetV2EXData ...
func GetV2EXData(num int) *hotso.HotData {
	data := internal.NewMongoDB().OnV2EXFindOne()
	var hotdata hotso.HotData
	if bytes, err := bson.MarshalJSON(data); err != nil {
		panic(err.Error())
	} else {
		bson.UnmarshalJSON(bytes, &hotdata)
	}
	var resultData []map[string]interface{}
	index := 0
	for _, v := range hotdata.Data.([]interface{}) {
		index++
		if num != 0 && index > num {
			break
		}
		resultData = append(resultData, v.(map[string]interface{}))
	}
	return &hotso.HotData{Type: hotdata.Type, Name: hotdata.Name, InTime: hotdata.InTime, Data: resultData}
}

//GetHotWordData ...
func GetHotWordData(c *gin.Context) {
	cli := internal.RedisCliPool().Get()
	defer cli.Close()
	var hottype = ""
	switch c.Param("hottype") {
	case "weibo":
		hottype = "weibo"
	case "baidu":
		hottype = "baidu"
	// case "zhihu":
	// 	hottype = "zhihu"
	default:
		c.JSON(http.StatusOK, gin.H{"code": -1, "message": "no data"})
	}
	// switch c.Param("data_type"){
	// case "json":

	// }
	num, _ := strconv.Atoi(c.Param("num"))
	key := internal.GetHotWordKey(hottype, time.Now().Year())
	if num <= 0 || num > 100 {
		num = 100
	}
	var result []map[string]interface{}
	args := redis.Args{}.Add(key).AddFlat([]string{"0", strconv.Itoa(num - 1), "WITHSCORES"})
	if reply, err := redis.Values(cli.Do("ZREVRANGE", args...)); err != nil {
		panic(err.Error())
	} else {
		var index = 0
		for i := 0; i < len(reply); i += 2 {
			index++
			result = append(result, map[string]interface{}{"rank": index, "word": string(reply[i].([]byte)), "score": string(reply[i+1].([]byte))})
		}
	}
	switch c.Param("data_type") {
	case "json":
		//c.JSON(http.StatusOK, result)
		ResponIndentJSON(c, http.StatusOK, result)
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "data format error",
		})
	}

}

//GetHotType ...
func GetHotType(c *gin.Context) {
	num, _ := strconv.Atoi(c.Param("num"))
	var data *hotso.HotData
	switch c.Param("hottype") {
	case "weibo":
		data = GetWeiBoData(num)
	case "baidu":
		data = GetBaiDuData(num)
	case "zhihu":
		data = GetZhiHuData(num)
	case "shuimu":
		data = GetShuiMuData(num)
	case "tianya":
		data = GetTianYaData(num)
	case "v2ex":
		data = GetV2EXData(num)
	default:
	}
	switch c.Param("data_type") {
	case "json":
		//c.JSON(http.StatusOK, data)
		ResponIndentJSON(c, http.StatusOK, data)
	// case "protobuf":
	default:
		c.JSON(http.StatusOK, gin.H{
			"hottype": c.Param("hottype"),
			"type":    c.Param("data_type"),
			"num":     c.Param("num"),
		})
	}
}

func main() {
	serviceCfg := config.GetConfig().Service
	router := gin.Default()
	v1 := router.Group("hotso/v1")
	{
		v1.GET("/hotso/:hottype/:data_type/:num", GetHotType) //  http://ip:port/weibo/json/10   获取微博热搜10条数据，并以json方式返回
		v1.GET("/hotword/:hottype/:data_type/:num", GetHotWordData)
	}
	addr := fmt.Sprintf("%s:%d", serviceCfg.IP, serviceCfg.Port)
	router.Run(addr)
}
