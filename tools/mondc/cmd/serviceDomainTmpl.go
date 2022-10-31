package cmd

import "text/template"

var domainEntityTemplate, _ = template.New("").Parse(`
package demo

import (
	"go.mongodb.org/mongo-driver/bson"
	mongodb "mond/wind/db/mongo"
	merr "mond/wind/err"
)

type Demo struct {
	Id        string ` + "`" + `bson:"id"` + "`" + `
	Status       int32 			 ` + "`" + `bson:"status" json:"status"` + "`" + `
	CreatedAt int64 			 ` + "`" + `bson:"createdAt" json:"createdAt"` + "`" + `
	UpdatedAt int64 			 ` + "`" + `bson:"updatedAt" json:"updatedAt"` + "`" + `
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

`)

var domainRepoTemplate, _ = template.New("").Parse(`
package demo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mongodb "mond/wind/db/mongo"
	merr "mond/wind/err"
)

type repository struct {
	mongoColl *mongodb.Collection
}

func newRepo(mongoColl *mongodb.Collection) *repository {
	repo := repository{
		mongoColl: mongoColl,
	}
	return &repo
}

func (m *repository) Find(ctx context.Context, condition interface{}, opts ...*options.FindOptions) ([]*Demo, error) {
	cursor, err := m.mongoColl.Find(ctx, condition, opts...)
	if err != nil {
		return nil, err
	}
	results := make([]*Demo, 0, 20)
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (m *repository) FindOne(ctx context.Context, condition interface{}, opts ...*options.FindOneOptions) (*Demo, error) {
	item := new(Demo)
	result, err := m.mongoColl.FindOne(ctx, condition, opts...)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, merr.FindOneResultNilError
	}
	//if result.Err() == mongo.ErrNoDocuments {
	//	return nil, merr.NoDocumentError
	//}
	err = result.Decode(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (m *repository) CountDocument(ctx context.Context, condition interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := m.mongoColl.CountDocuments(ctx, condition, opts...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *repository) InsertOne(ctx context.Context, item *Demo) (*Demo, error) {
	result, err := m.mongoColl.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}
	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, merr.TransObjectIdError
	}
	return item, nil
}

func (m *repository) InsertMany(ctx context.Context, items []*Demo) error {
	documents := make([]interface{}, 0, len(items))
	for _, value := range items {
		documents = append(documents, value)
	}
	_, err := m.mongoColl.InsertMany(ctx, documents)
	if err != nil {
		return err
	}
	return nil
}
func (m *repository) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	r, err := m.mongoColl.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}
	return r.ModifiedCount, nil
}

func (m *repository) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	r, err := m.mongoColl.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}
	return r.ModifiedCount, err
}

func (m *repository) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	cursor, err := m.mongoColl.Aggregate(ctx, pipeline, opts...)
	return cursor, err
}


`)

var domainServiceTemplate, _ = template.New("").Parse(`
package demo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	mredis "mond/wind/cache/redis"
	mongodb "mond/wind/db/mongo"
)

type Service struct {
	repo     *repository
	redisCli *mredis.Client
}

func NewService(mongoColl *mongodb.Collection, redisCli *mredis.Client) *Service {
	svc := Service{}
	svc.repo = newRepo(mongoColl)
	svc.redisCli = redisCli
	return &svc
}

func (m *Service) FindOneByCondition(ctx context.Context, filter *option, opts ...*options.FindOneOptions) (*Demo, error) {
	if filter.err != nil {
		return nil, filter.err
	}
	item, err := m.repo.FindOne(ctx, filter.pm.GetFilter(), opts...)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (m *Service) FindByCondition(ctx context.Context, filter *option, opts ...*options.FindOptions) ([]*Demo, error) {
	if filter.err != nil {
		return nil, filter.err
	}
	items, err := m.repo.Find(ctx, filter.pm.GetFilter(), opts...)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (m *Service) CreateOne(ctx context.Context, document *Demo) error {
	_, err := m.repo.InsertOne(ctx, document)
	return err
}

func (m *Service) UpdateOneById(ctx context.Context, id string, update *option, opts ...*options.UpdateOptions) error {
	if update.err != nil {
		return update.err
	}
	_, err := m.repo.UpdateOne(ctx, bson.M{"_id": id}, update.pm.GetChangeInfo(), opts...)
	if err != nil {
		return err
	}
	return err
}

func (m *Service) UpdateByCondition(ctx context.Context, update *option, opts ...*options.UpdateOptions) (int64, error) {
	if update.err != nil {
		return 0, update.err
	}
	count, err := m.repo.UpdateMany(ctx, update.pm.GetFilter(), update.pm.GetChangeInfo(), opts...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *Service) CountByCondition(ctx context.Context, filter *option, opts ...*options.CountOptions) (int64, error) {
	if filter.err != nil {
		return 0, filter.err
	}
	count, err := m.repo.CountDocument(ctx, filter.pm.GetFilter(), opts...)
	return count, err
}


`)
