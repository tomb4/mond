package cmd

import "text/template"

var handlerHandlerTemplate, _ = template.New("").Parse(`
package handler

import (
	"context"
	"mond/wind/logger"
	"mond/service/{{.FolderPath}}/app"
	"mond/service/{{.FolderPath}}/proto"
)

type {{.AppId}}Service struct {
	_log logger.Logger
	app *app.App
}

func (m *{{.AppId}}Service) Ping(ctx context.Context, n *{{.AppId}}.PingReq) (*{{.AppId}}.PingResp, error) {
	panic("implement me")
}

func New{{.AppId}}Service() *{{.AppId}}Service {
	return &{{.AppId}}Service{_log: logger.GetLogger()}
}


`)

var handlerTestTemplate, _ = template.New("").Parse(`
package handler

import (
	"context"
	"fmt"
	"mond/wind"
	"mond/service/{{.FolderPath}}/proto"
	"os"
	"testing"
)

var (
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	h := NewHook()
	wind.InitFrame(h, wind.WithTestMode())
	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test{{.AppId}}Service_Ping(t *testing.T) {
	req := &{{.AppId}}.PingReq{}
	resp, err := _handler.Ping(ctx, req)
	t.Error(err)
	fmt.Println(resp)
}

`)
var handlerHookTemplate, _ = template.New("").Parse(`
package handler

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	hook2 "mond/wind/hook"
	"mond/wind/logger"
	dynamicConfig "mond/service/{{.FolderPath}}/infra/config"
	"mond/service/{{.FolderPath}}/proto"
)

type hook struct {
	h *{{.AppId}}Service
}

var (
	_handler *{{.AppId}}Service
)

func NewHook() hook2.FrameStartHook {
	_handler = New{{.AppId}}Service()
	return &hook{
		h: _handler,
	}
}

//框架基础组件初始化完成后，加载应用资源
func (h *hook) ResourceInitHook(ctx context.Context) error {
	err := h.h.ResourceInit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (h *hook) GrpcStartHook(s *grpc.Server) error {
	{{.AppId}}.Register{{.AppId}}ServiceServer(s, {{.AppId}}.New{{.AppId}}ServiceServerImpl(h.h))

	return nil
}

func (h *hook) ConfigChangeHook(conf map[string]interface{}) {
	c := &dynamicConfig.DynamicConfig{}
	bytes, err := json.Marshal(conf)
	if err != nil {
		logger.GetLogger().Error(context.TODO(), "ConfigChangeHook失败", zap.Any("err", err))
		return
	}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		logger.GetLogger().Error(context.TODO(), "ConfigChangeHook失败2", zap.Any("err", err))
		return
	}
	dynamicConfig.SetDynamicConfig(c)
}
func (h *hook) AppStopHook() {

}

`)
var handlerResourceTemplate, _ = template.New("").Parse(`
package handler

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"mond/wind"
	"mond/service/{{.FolderPath}}/app"
	"mond/service/{{.FolderPath}}/domain/demo"
	"sync"
	"time"
)

var (
	resourceLock sync.Mutex
	resourceInit bool
)

func (m *{{.AppId}}Service) ResourceInit(ctx context.Context) error {
	resourceLock.Lock()
	defer resourceLock.Unlock()
	if resourceInit {
		panic("ResourceInit已经执行过了")
	}
	rdb, err := wind.GetRedisDbManager().GetClient("master")
	if err != nil {
		return err
	}

	dbm := wind.GetMongoDbManager()
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	demoColl, err := dbm.GetCollection("master.meta.demo")
	if err != nil {
		return err
	}
	index := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "id", Value: bsonx.Int32(1)},
			},
			Options: options.Index().SetBackground(true).SetUnique(true),
		},
	}
	demoColl.Indexes().CreateMany(ctx, index, opts)

	demoSvc := demo.NewService(demoColl, rdb)

	m.app = app.NewApp(rdb, demoSvc)

	resourceInit = true
	return nil
}


`)
