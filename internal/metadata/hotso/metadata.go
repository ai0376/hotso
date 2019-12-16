package hotso

//HotData ...
type HotData struct {
	Type   int         `json:"type"`
	Name   string      `json:"name"`
	InTime int64       `json:"intime"`
	Data   interface{} `json:"data"`
}

//enum
const (
	WEIBO  = iota //微博
	BAIDU         //百度
	ZHIHU         //知乎
	SHUIMU        //水木
	TIANYA        //天涯
	V2EX          //v2ex
)

//HotSoType ...
var HotSoType = map[int]string{
	WEIBO:  "WeiBo",
	BAIDU:  "BaiDu",
	ZHIHU:  "ZhiHu",
	SHUIMU: "ShuiMu",
	TIANYA: "TianYa",
	V2EX:   "V2EX",
}

//HotItem ...
type HotItem struct {
	Reading string `json:"reading"`
	State   string `json:"state"`
	Title   string `json:"title"`
	Top     string `json:"top"`
	URL     string `json:"url"`
}
