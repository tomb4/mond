package async

import (
	"context"
	"fmt"
	"mond/wind/utils"
	"reflect"
	"testing"
)

func TestAsync_Register(t *testing.T) {
	a := async{
		funcInputTypeMap: map[string]reflect.Type{},
		funcMap:          map[string]reflect.Value{},
	}
	app := &App{}
	a.Register(app)
	err := app.SelectUserAsync(context.TODO(), &Dto{Id: 123})
	utils.MustNil(err)
}

type App struct {
	SelectUserAsync func(context.Context, *Dto) error `async:"true"`
}
type Dto struct {
	Id int32
}

func (m *App) SelectUser(ctx context.Context, dto *Dto) error {
	fmt.Println("SelectUser", dto)
	return nil
}
