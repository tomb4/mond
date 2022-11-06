package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"mond/wind/config"
	"mond/wind/env"
	merr "mond/wind/err"
	"mond/wind/utils/endpoint"
	"time"
)

type Client struct {
	cli *mongo.Client
}

func (m *dbManager) GetClient(instance string) (*Client, error) {
	if env.GetAppState() != env.Starting {
		panic("mongo GetClient只能在ResourceInit中调用")
	}
	m.clientLock.Lock()
	defer m.clientLock.Unlock()
	
	if m.clientMap[instance] != nil {
		return m.clientMap[instance], nil
	}
	opt, err := config.GetMongoClientOption(instance)
	if err != nil {
		return nil, err
	}
	if opt.Uri == "" {
		return nil, errors.New("mongo uri cannot empty")
	}
	opts := options.Client().ApplyURI(opt.Uri)
	if opt.ReadPreference != "" {
		var _readpref *readpref.ReadPref
		var readprefOption readpref.Option
		if opt.MaxStaleness != 0 {
			readprefOption = readpref.WithMaxStaleness(time.Millisecond * time.Duration(opt.MaxStaleness))
		}
		switch opt.ReadPreference {
		case "Primary":
			_readpref = readpref.Primary()
		case "PrimaryPreferred":
			_readpref = readpref.PrimaryPreferred(readprefOption)
		case "SecondaryPreferred":
			_readpref = readpref.SecondaryPreferred(readprefOption)
		case "Secondary":
			_readpref = readpref.Secondary(readprefOption)
		}
		opts.SetReadPreference(_readpref)
	}
	//TODO:  还有其它参数的使用可以接入
	c, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}
	err = c.Connect(context.TODO())
	if err != nil {
		return nil, err
	}
	err = c.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, err
	}
	m.clientMap[instance] = &Client{cli: c}
	return m.clientMap[instance], nil
}

func (c *Client) UseSession(ctx context.Context, fn func(sessionCtx SessionCtx) error) error {
	select {
	case <-ctx.Done():
		return merr.MongoContextTimeoutErr
	default:
	}
	return c.cli.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		return fn(NewSessionCtx(sessionContext))
	})
}

type SessionCtx struct {
	mongo.SessionContext
	withTransactionEndpoint endpoint.Endpoint
}

func NewSessionCtx(sessionContext mongo.SessionContext) SessionCtx {
	s := SessionCtx{}
	s.SessionContext = sessionContext
	
	s.withTransactionEndpoint = s.makeWithTransactionEndpoint()
	
	return s
}

func (s SessionCtx) makeWithTransactionEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		_mongoReq := req.(*mongoReq)
		opts := _mongoReq.opts
		_opts := make([]*options.TransactionOptions, 0, len(opts))
		for _, v := range opts {
			_opts = append(_opts, v.(*options.TransactionOptions))
		}
		fun := _mongoReq.withTFunc.(func(sessCtx mongo.SessionContext) (interface{}, error))
		return s.WithTransaction(ctx, fun, _opts...)
	}
}

func (s SessionCtx) Transaction(ctx context.Context, fn func(sessionContext mongo.SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions) error {
	// TODO ctx inject
	_opts := make([]interface{}, 0, len(opts))
	_opts = append(_opts, options.Transaction().
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetReadConcern(readconcern.Snapshot()).SetReadPreference(readpref.Primary()))
	for _, v := range opts {
		_opts = append(_opts, v)
	}
	req := &mongoReq{
		opts:      _opts,
		withTFunc: fn,
	}
	
	_, err := s.withTransactionEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
