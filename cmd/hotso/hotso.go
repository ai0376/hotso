package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/mjrao/hotso/common"
	"github.com/mjrao/hotso/config"
	"github.com/mjrao/hotso/internal"
	"github.com/mjrao/hotso/internal/cloud"
	"github.com/mjrao/hotso/internal/metadata/hotso"
)

//Spider ...
type Spider struct {
	Type int
}

var wg *sync.WaitGroup
var userAgent = "Chrome: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"

//OnWeiBo ...
func (s *Spider) OnWeiBo() []map[string]interface{} {
	url := "https://s.weibo.com/top/summary"

	var allData []map[string]interface{}

	c := colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(userAgent))
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnHTML("#pl_top_realtimehot > table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tbody > tr", func(i int, ex *colly.HTMLElement) {
			top := ex.ChildText("td.td-01.ranktop")
			title := ex.ChildText("td.td-02 > a")
			reading := ex.ChildText("td.td-02 > span")
			state := ex.ChildText("td.td-03 > i")
			var url = ""
			if state == "荐" { //广告数据
				url = ex.ChildAttr("td.td-02 > a", "href_to")
			} else {
				url = ex.ChildAttr("td.td-02 > a", "href")
			}
			allData = append(allData, map[string]interface{}{"top": top, "title": title, "reading": reading, "url": "https://s.weibo.com" + url, "state": state})
		})
	})
	c.Visit(url)
	return allData
}

//OnBaiDu 实时热点
func (s *Spider) OnBaiDu() []map[string]interface{} {
	url := "http://top.baidu.com/buzz?b=1&c=513&fr=topbuzz_b341_c513"
	var allData []map[string]interface{}

	c := colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(userAgent))
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnHTML("#main > div.mainBody > div > table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tbody > tr", func(i int, ex *colly.HTMLElement) {
			top := ex.ChildText("td.first > span")
			if top != "" {
				title := ex.ChildText("td.keyword > a.list-title")
				reading := ex.ChildText("td.last > span")
				url := ex.ChildAttr("td.keyword > a.list-title", "href")
				state := "" //ex.ChildText("td.td-03 > i")
				allData = append(allData, map[string]interface{}{"top": top, "title": common.GBK2UTF8(title), "reading": reading, "url": url, "state": state})
			}
		})
	})
	c.Visit(url)
	return allData
}

//OnZhiHu 实时热点
func (s *Spider) OnZhiHu() []map[string]interface{} {

	//ZhiHuOnline ...
	type ZhiHuOnline struct {
		Cookie    string `json:"cookie"`
		UserAgent string `json:"user_agent"`
	}

	var allData []map[string]interface{}
	var success = true

	var zhihu ZhiHuOnline
	if webdavCli, err := cloud.Dial(config.GetConfig().WebDav.Host, config.GetConfig().WebDav.User, config.GetConfig().WebDav.Password); err != nil {
		fmt.Println("zhihu webdav dial error")
		success = false
	} else {
		remoteDir := strings.Replace(config.GetConfig().WebDav.RemoteDir, "\\", "/", -1)
		if remoteDir[len(remoteDir)-1:] != "/" {
			remoteDir = remoteDir + "/"
		}
		if body, err := webdavCli.Download(remoteDir + "zhihu.json"); err != nil {
			fmt.Println("zhihu webdav download error")
			success = false
		} else {
			json.Unmarshal(body, &zhihu)
		}
	}
	if success != true {
		return allData
	}

	c := colly.NewCollector(colly.UserAgent(zhihu.UserAgent), colly.MaxDepth(1))
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", zhihu.Cookie)
	})
	c.OnHTML("#TopstoryContent > div > div > div.HotList-list", func(e *colly.HTMLElement) {
		e.ForEach("div.HotList-list > section.HotItem", func(i int, ex *colly.HTMLElement) {
			top := ex.ChildText("div.HotItem-index > div.HotItem-rank")
			title := ex.ChildText("div.HotItem-content > a > h2.HotItem-title")
			hotread := ex.ChildText("div.HotItem-content > div.HotItem-metrics")
			var reading = 0
			var err error
			ss := strings.Fields(hotread)
			if len(ss) >= 2 {
				if index := strings.Index(hotread, "万"); index == -1 {
					if reading, err = strconv.Atoi(ss[0]); err != nil {
						fmt.Println("zhihu hotnum error")
					}
				} else {
					if reading, err = strconv.Atoi(ss[0]); err != nil {
						fmt.Println("zhihu hotnum error")
					} else {
						reading = reading * 10000
					}
				}
			}
			url := ex.ChildAttr("div.HotItem-content > a ", "href")
			state := ex.ChildText("div.HotItem-index > div.HotItem-label")

			allData = append(allData, map[string]interface{}{"top": top, "title": title, "reading": reading, "url": "https://s.weibo.com" + url, "state": state})
		})
	})
	c.Visit("http://www.zhihu.com/hot")

	return allData
}

//ProduceData ...
func ProduceData(s *Spider) {
	defer wg.Done()
	reflectValue := reflect.ValueOf(s)
	methodValue := reflectValue.MethodByName("On" + hotso.HotSoType[s.Type])
	methodFunc := methodValue.Call(nil)
	originData := methodFunc[0].Interface().([]map[string]interface{}) //数据
	now := time.Now().Unix()
	if len(originData) > 0 {
		switch s.Type {
		case hotso.WEIBO:
			internal.NewMongoDB().OnWeiBoInsert(&hotso.HotData{Type: s.Type, Name: hotso.HotSoType[s.Type], InTime: now, Data: originData})
		case hotso.BAIDU:
			internal.NewMongoDB().OnBaiDuInsert(&hotso.HotData{Type: s.Type, Name: hotso.HotSoType[s.Type], InTime: now, Data: originData})
		case hotso.ZHIHU:
			internal.NewMongoDB().OnZhiHuInsert(&hotso.HotData{Type: s.Type, Name: hotso.HotSoType[s.Type], InTime: now, Data: originData})
		}
	} else {
		fmt.Println("originData nil")
	}
}

func main() {
	wg = &sync.WaitGroup{}
	wg.Add(len(hotso.HotSoType))
	for k, _ := range hotso.HotSoType {
		s := &Spider{Type: k}
		go ProduceData(s)
	}
	wg.Wait()
}
