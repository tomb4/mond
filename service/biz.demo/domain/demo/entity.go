package demo

import (
	"go.mongodb.org/mongo-driver/bson"
	mongodb "mond/wind/db/mongo"
	merr "mond/wind/err"
)

type Demo struct {
	Id        string `bson:"id"`
	Status    int32  `bson:"status" json:"status"`
	CreatedAt int64  `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64  `bson:"updatedAt" json:"updatedAt"`
}

type option struct {
	pm  *mongodb.PatchMode
	err error
}

func Option() *option {
	return &option{pm: mongodb.NewPatchMode()}
}

func (m *option) FilterId(v string) *option {
	if v == "" {
		m.err = merr.DomainOptionError.WithMsg("FilterId")
		return m
	}
	m.pm.Filter("id", v)
	return m
}

func (m *option) SetStatus(v int32) *option {
	if v <= 0 || v > 2 {
		m.err = merr.DomainOptionError.WithMsg("SetStatus")
		return m
	}
	m.pm.Set("status", v)
	return m
}

func (m *option) SetUpdatedAt(v int64) *option {
	if v <= 0 {
		m.err = merr.DomainOptionError.WithMsg("SetUpdatedAt")
		return m
	}
	m.pm.Set("updatedAt", v)
	return m
}

func (m *option) FilterGtCreatedAt(v int64) *option {
	if v <= 0 {
		m.err = merr.DomainOptionError.WithMsg("FilterGtCreatedAt")
		return m
	}
	m.pm.Filter("createdAt", bson.M{"$gt": v})
	return m
}
