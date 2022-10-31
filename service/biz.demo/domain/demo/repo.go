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
