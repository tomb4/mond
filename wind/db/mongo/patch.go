package mongodb

import "go.mongodb.org/mongo-driver/bson"

type PatchMode struct {
	filter   bson.M
	chInfo   bson.M
	document bson.M
	cacheKey string
	//opt      interface{}
}

func NewPatchMode() *PatchMode {
	return &PatchMode{}
}

func (m *PatchMode) SetCKey(key string) {
	m.cacheKey = key
}

func (m *PatchMode) GetCKey() string {
	return m.cacheKey
}

func (m *PatchMode) GetFilter() *bson.M {
	return &m.filter
}

func (m *PatchMode) GetChangeInfo() *bson.M {
	return &m.chInfo
}

func (m *PatchMode) GetDocument() *bson.M {
	return &m.document
}

//func (m *PatchMode) GetOpt() interface{} {
//	return m.opt
//}

func (m *PatchMode) Clear() {
	m.filter = nil
	m.chInfo = nil
	//m.opt = nil
	m.document = nil
}

func (m *PatchMode) _filter(field string, value interface{}) {
	if m.filter == nil {
		m.filter = bson.M{}
	}
	m.filter[field] = value
}

func (m *PatchMode) _update(_type, field string, value interface{}) {
	if m.chInfo == nil {
		m.chInfo = bson.M{}
	}
	if set, ok := m.chInfo[_type]; ok {
		if setV, ok := set.(bson.M); ok {
			setV[field] = value
		}
	} else {
		setV := bson.M{
			field: value,
		}
		m.chInfo[_type] = setV
	}
	m._document(field, value)
}

func (m *PatchMode) _document(field string, value interface{}) {
	if m.document == nil {
		m.document = bson.M{}
	}
	m.document[field] = value
}

func (m *PatchMode) Filter(field string, value interface{}) {
	m._filter(field, value)
}

func (m *PatchMode) Set(field string, value interface{}) {
	m._update("$set", field, value)
}

func (m *PatchMode) Inc(field string, num interface{}) {
	m._update("$inc", field, num)
}

func (m *PatchMode) Pull(field string, value interface{}) {
	m._update("$pull", field, value)
}

func (m *PatchMode) Push(field string, value interface{}) {
	m._update("$push", field, value)
}

func (m *PatchMode) AddToSet(field string, value interface{}) {
	m._update("$addToSet", field, value)
}

func (m *PatchMode) SetOnInsert(field string, value interface{}) {
	m._update("$setOnInsert", field, value)
}

func (m *PatchMode) SetMulti(kw bson.M) {
	if m.chInfo == nil {
		m.chInfo = bson.M{}
	}
	if set, ok := m.chInfo["$set"]; ok {
		if setV, ok := set.(bson.M); ok {
			for field, value := range kw {
				setV[field] = value
			}
		}
	} else {
		m.chInfo["$set"] = kw
	}
}

//func (m *PatchMode) Opt(opt interface{}) {
//	m.opt = opt
//}

