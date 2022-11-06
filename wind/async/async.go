package async

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mond/wind/env"
	"mond/wind/mq/rabbit"
	"mond/wind/utils"
	"mond/wind/utils/constant"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

var (
	nilError = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
)

type Async interface {
	Register(svc interface{})
}

type async struct {
	rabbitClient     rabbit.Client
	funcMap          map[string]reflect.Value
	funcInputTypeMap map[string]reflect.Type
	lock             sync.Mutex
}

func (m *async) Register(svc interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	valueOf := reflect.ValueOf(svc)
	valueOfElem := valueOf.Elem()
	if valueOf.Kind() != reflect.Ptr || valueOfElem.Kind() != reflect.Struct {
		panic("注册的必须是结构体的指针")
	}
	typeOf := reflect.TypeOf(svc)
	typeOfElem := typeOf.Elem()
	//寻找有tag的异步函数
	for i := 0; i < typeOfElem.NumField(); i++ {
		field := typeOfElem.Field(i)
		tag := field.Tag
		if tag.Get("async") != "true" {
			continue
		}
		if valueOfElem.Field(i).Kind() != reflect.Func {
			panic("异步函数必须是函数")
		}
		name := field.Name
		if !unicode.IsUpper([]rune(name)[0]) {
			panic("异步函数首字母必须大写")
		}
		if !strings.Contains(name, "Async") || strings.Index(name, "Async") != len(name)-5 {
			panic("异步函数必须以Async结尾")
		}
		syncMethod := strings.ReplaceAll(name, "Async", "")
		if m.funcMap[syncMethod].IsValid() {
			panic("不同结构体中有方法名相同的异步任务")
		}
		
		method := valueOf.MethodByName(syncMethod)
		if !method.IsValid() || method.IsZero() || method.IsNil() {
			panic(fmt.Sprintf("与异步任务对应的函数%s不存在", syncMethod))
		}
		if field.Type != method.Type() {
			panic("异步函数与同步函数不匹配")
		}
		if method.Type().NumOut() != 1 {
			panic("异步函数返回参数只能有error")
		}
		if method.Type().NumIn() != 2 {
			panic("异步函数入参必须只能有2位")
		}
		//if !(method.Type().In(1).Kind() == reflect.Struct || (method.Type().In(1).Kind() == reflect.Ptr && method.Type().In(1).Elem().Kind() == reflect.Struct)) {
		if !(method.Type().In(1).Kind() == reflect.Ptr && method.Type().In(1).Elem().Kind() == reflect.Struct) {
			//panic("异步函数入参第二位必须为struct或者struct的指针")
			panic("异步函数入参第二位必须为struct的指针")
		}
		m.funcMap[name] = method
		m.funcInputTypeMap[name] = method.Type().In(1)
		
		valueOfElem.Field(i).Set(reflect.MakeFunc(method.Type(), func(args []reflect.Value) (results []reflect.Value) {
			ctx, ok := args[0].Interface().(context.Context)
			if !ok {
				results = append(results, reflect.ValueOf(errors.New("第一个入参必须是context")))
				return
			}
			if args[1].IsNil() {
				results = append(results, reflect.ValueOf(errors.New("参数不能为nil")))
				return
			}
			input := args[1].Interface()
			bytes, err := json.Marshal(input)
			if err != nil {
				results = append(results, reflect.ValueOf(err))
				return
			}
			msg := &AsyncMessage{Bytes: bytes, Method: name}
			ctx = context.WithValue(ctx, constant.AsyncMethodCtxKey, name)
			err = m.rabbitClient.Publish(ctx, msg)
			//下面是为了本地单元测试的
			//err = m.execute(ctx, &consumeMsg{
			//	Delivery: &amqp091.Delivery{
			//		Body: msg.Body(),
			//	},
			//})
			if err != nil {
				results = append(results, reflect.ValueOf(err))
				return
			}
			results = append(results, nilError)
			return
		}))
		
	}
	err := m.rabbitClient.InitProducer()
	utils.MustNil(err)
	err = m.rabbitClient.CreateExchange(constant.TopicName, "topic", true, false)
	utils.MustNil(err)
	err = m.rabbitClient.CreateQueue(genQueue(), true, false)
	utils.MustNil(err)
	err = m.rabbitClient.QueueBind(genQueue(), genQueue(), constant.TopicName)
	utils.MustNil(err)
	err = m.rabbitClient.Consume(genQueue(), false, m.execute)
	if err != nil {
		panic(err)
	}
}

func genQueue() string {
	return fmt.Sprintf("Async.%s", env.GetAppId())
}

//type consumeMsg struct {
//	*amqp091.Delivery
//	ack bool
//}
//
//func (m *consumeMsg) Exchange() string {
//	return m.Delivery.Exchange
//}
//
//func (m *consumeMsg) RoutingKey() string {
//	return m.Delivery.RoutingKey
//}
//
//func (m *consumeMsg) Body() []byte {
//	return m.Delivery.Body
//}
//
//func (m *consumeMsg) Nack(multiple, requeue bool) error {
//	if m.ack {
//		return nil
//	}
//	m.ack = true
//	return m.Delivery.Nack(multiple, requeue)
//}
//
//func (m *consumeMsg) Ack(multiple bool) error {
//	if m.ack {
//		return nil
//	}
//	m.ack = true
//	return m.Delivery.Ack(multiple)
//}

func (m *async) execute(ctx context.Context, msg rabbit.Message) error {
	defer msg.Ack(false)
	asyncMsg := AsyncMessage{}
	err := json.Unmarshal(msg.Body(), &asyncMsg)
	if err != nil {
		return err
	}
	method := m.funcMap[asyncMsg.Method]
	if !method.IsValid() {
		msg.Nack(false, true)
		return errors.New("method未找到")
	}
	
	inputPtr := reflect.New(m.funcInputTypeMap[asyncMsg.Method].Elem())
	err = json.Unmarshal(asyncMsg.Bytes, inputPtr.Interface())
	if err != nil {
		return err
	}
	results := method.Call([]reflect.Value{reflect.ValueOf(ctx), inputPtr})
	if results[0].IsValid() && !results[0].IsNil() {
		return results[0].Interface().(error)
	}
	return nil
}

type AsyncMessage struct {
	Bytes  []byte `json:"bytes"`
	Method string `json:"method"`
}

func (m *AsyncMessage) Exchange() string {
	return constant.TopicName
}

func (m *AsyncMessage) RoutingKey() string {
	return genQueue()
}

func (m *AsyncMessage) Body() []byte {
	bytes, _ := json.Marshal(m)
	
	return bytes
}

func GetRabbitAsyncClient(rabbitClient rabbit.Client) (Async, error) {
	if env.GetAppState() != env.Starting {
		panic("GetAsyncClient只能在ResourceInit中调用")
	}
	_async := async{
		rabbitClient:     rabbitClient,
		funcMap:          map[string]reflect.Value{},
		funcInputTypeMap: map[string]reflect.Type{},
	}
	return &_async, nil
}
