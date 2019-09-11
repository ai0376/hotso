package hotso

//HotData ...
type HotData struct {
	Type   int
	Name   string
	InTime int64
	Data   interface{}
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
