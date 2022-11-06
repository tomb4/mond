package mredis

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
    "go.mongodb.org/mongo-driver/mongo"
    "testing"
    "time"
)

var (
    ctx = context.TODO()
)

type Cache struct {
    Name string
}

func TestNewRedisCache(t *testing.T) {
    rdbCli := redis.NewClient(&redis.Options{
        Addr: "dev-api.neoclub.cn:6379",
        DB:   7, // use default DB

    })
    //var num int64
    c := NewRedisCache(&Client{Client: rdbCli}, mongo.ErrNoDocuments)
    //t.Run("not found", func(t *testing.T) {
    //    err := c.Take(ctx, "demo not found", &num, func(v interface{}) error {
    //        return mongo.ErrNoDocuments
    //    })
    //    if err != mongo.ErrNoDocuments {
    //        t.Error(err)
    //    }
    //})
    //t.Run("other err", func(t *testing.T) {
    //    err := c.Take(ctx, "demo other err", &num, func(v interface{}) error {
    //        return errors.New("other error")
    //    })
    //    if err != errors.New("other error") {
    //        t.Error(err)
    //    }
    //})
    var ca []*Cache = make([]*Cache, 0)
    t.Run("normal", func(t *testing.T) {
        err := c.Take(ctx, fmt.Sprintf("demo %d", time.Now().Unix()), &ca, func(v interface{}) error {
            //cc := &Cache{}
            //cc.Name = fmt.Sprintf("c1 %d", time.Now().Unix())
            //bytes, _ := json.Marshal(cc)
            //json.Unmarshal(bytes, v)
            fmt.Println(v)
            return nil
        })
        t.Error(err)
    })
    fmt.Println(ca)
}
