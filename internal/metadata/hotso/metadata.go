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
)

//HotSoType ...
var HotSoType = map[int]string{
	WEIBO:  "WeiBo",
	BAIDU:  "BaiDu",
	ZHIHU:  "ZhiHu",
	SHUIMU: "ShuiMu",
	TIANYA: "TianYa",
}
