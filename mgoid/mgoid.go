package mgoid

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// IncID 增量ID
type IncID struct {
	session *mgo.Session
	db      *mgo.Database
	c       *mgo.Collection
}

// NewIncID 创建IncID的实例
func NewIncID(mgoURL, dbName string) (*IncID, error) {
	return NewIncIDWithC(mgoURL, dbName, "counters")
}

// NewIncIDWithC 创建IncID的实例(提供存储的集合名称)
func NewIncIDWithC(mgoURL, dbName, cName string) (*IncID, error) {
	session, err := mgo.Dial(mgoURL)
	if err != nil {
		return nil, err
	}
	db := session.DB(dbName)
	return &IncID{
		session: session,
		db:      db,
		c:       db.C(cName),
	}, nil
}

// GenerateID 生成自增ID
func (i *IncID) GenerateID(name string) (id int64, err error) {
	var result incResult
	_, err = i.c.Find(bson.M{"_id": name}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		ReturnNew: true,
		Upsert:    true,
	}, &result)
	if err != nil {
		return
	}
	id = result.Seq
	return
}

// Close 关闭会话
func (i *IncID) Close() {
	i.session.Close()
}

type incResult struct {
	Seq int64 `bson:"seq"`
}
