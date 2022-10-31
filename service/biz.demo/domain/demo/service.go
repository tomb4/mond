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
