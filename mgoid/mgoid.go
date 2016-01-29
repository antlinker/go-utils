package mgoid

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	_GCollectionName = "counters"
)

// IncID 增量ID
type IncID struct {
	session *mgo.Session
	db      *mgo.Database
	c       *mgo.Collection
}

// NewIncIDWithURL 创建IncID的实例
func NewIncIDWithURL(url, dbName string) (*IncID, error) {
	return NewIncIDWithURLC(url, dbName, _GCollectionName)
}

// NewIncIDWithURLC 创建IncID的实例
func NewIncIDWithURLC(url, dbName, cName string) (*IncID, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	db := session.DB(dbName)
	return NewIncIDWithSession(session, db, cName), nil
}

// NewIncID 创建IncID的实例
func NewIncID(db *mgo.Database) *IncID {
	return NewIncIDWithDB(db, _GCollectionName)
}

// NewIncIDWithDB 创建IncID的实例
func NewIncIDWithDB(db *mgo.Database, cName string) *IncID {
	return NewIncIDWithSession(db.Session.Clone(), db, cName)
}

// NewIncIDWithSession 创建IncID的实例
func NewIncIDWithSession(session *mgo.Session, db *mgo.Database, cName string) *IncID {
	return &IncID{
		session: session,
		db:      db,
		c:       db.C(cName),
	}
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

func (i *IncID) Close() {
	i.session.Close()
}

type incResult struct {
	Seq int64 `bson:"seq"`
}
