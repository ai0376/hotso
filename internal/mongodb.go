package internal

import (
	"github.com/mjrao/hotso/config"
	"github.com/mjrao/hotso/internal/metadata/hotso"
	"gopkg.in/mgo.v2"
)

//MongoDB ...
type MongoDB struct {
}

var session *mgo.Session

func init() {
	if session == nil {
		if s, err := mgo.Dial(config.GetConfig().MongoDB.Host); err != nil {
			panic(err.Error())
		} else {
			session = s
		}
	}
}

//NewMongoDB ...
func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

//OnWeiBoInsert ...
func (m *MongoDB) OnWeiBoInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("weibo").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnWeiBoFindOne ...
func (m *MongoDB) OnWeiBoFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("weibo").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//OnBaiDuInsert ...
func (m *MongoDB) OnBaiDuInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("baidu").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnBaiDuFindOne ...
func (m *MongoDB) OnBaiDuFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("baidu").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//OnZhiHuInsert ...
func (m *MongoDB) OnZhiHuInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("zhihu").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnZhiHuFindOne ...
func (m *MongoDB) OnZhiHuFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("zhihu").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//OnShuiMuInsert ...
func (m *MongoDB) OnShuiMuInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("shuimu").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnShuiMuFindOne ...
func (m *MongoDB) OnShuiMuFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("shuimu").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//OnTianYaInsert ...
func (m *MongoDB) OnTianYaInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("tianya").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnTianYaFindOne ...
func (m *MongoDB) OnTianYaFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("tianya").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//OnV2EXInsert ...
func (m *MongoDB) OnV2EXInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("v2ex").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnV2EXFindOne ...
func (m *MongoDB) OnV2EXFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("v2ex").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//---------hottop-----------
//OnLoadData ...
func (m *MongoDB) OnLoadData(dataType int, begin int64, end int64) *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	// collection := strings.ToLower(hotso.HotSoType[dataType])
	// if err:= s.DB("hotso").C(collection).Find()
	var data hotso.HotData
	return &data
}
