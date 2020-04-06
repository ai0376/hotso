package internal

import (
	"fmt"
	"strings"

	"github.com/mjrao/hotso/config"
	"github.com/mjrao/hotso/internal/metadata/hotso"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

//OnInsertDataByType ...
func (m *MongoDB) OnInsertDataByType(dataType int, data interface{}) {
	s := session.Copy()
	defer s.Close()
	if _, ok := hotso.HotSoType[dataType]; !ok {
		return
	}
	col := strings.ToLower(hotso.HotSoType[dataType])
	if err := s.DB("hotso").C(col).Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnFindOneDataByType ...
func (m *MongoDB) OnFindOneDataByType(dataType int) *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	if _, ok := hotso.HotSoType[dataType]; !ok {
		return nil
	}
	col := strings.ToLower(hotso.HotSoType[dataType])
	var data hotso.HotData
	if err := s.DB("hotso").C(col).Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//OnQueryData ...
func (m *MongoDB) OnQueryData(dataType int, start int64, end int64) *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	if _, ok := hotso.HotSoType[dataType]; !ok {
		return nil
	}

	collection := strings.ToLower(hotso.HotSoType[dataType])
	var data hotso.HotData
	if err := s.DB("hotso").C(collection).Find(bson.M{"intime": bson.M{"$gt": start, "$lte": end}}).Sort("-intime").Limit(1).One(&data); err != nil {
		fmt.Println(err.Error())
	}
	return &data
}

//---------hottop-----------
//OnLoadData ...
func (m *MongoDB) OnLoadData(dataType int, begin int64, end int64) []hotso.HotData {
	s := session.Copy()
	defer s.Close()
	collection := strings.ToLower(hotso.HotSoType[dataType])
	var datas []hotso.HotData
	if err := s.DB("hotso").C(collection).Find(bson.M{"intime": bson.M{"$gt": begin, "$lte": end}}).All(&datas); err != nil {
		fmt.Println(err.Error())
	}
	return datas
}
