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
	WEIBO = iota
	BAIDU
	ZHIHU
)

//HotSoType ...
var HotSoType = map[int]string{
	WEIBO: "WeiBo",
	BAIDU: "BaiDu",
	ZHIHU: "ZhiHu",
}
