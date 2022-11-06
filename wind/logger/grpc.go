package logger

import (
	"context"
	"encoding/json"
	"fmt"
	merr "mond/wind/err"
	"mond/wind/utils"
	byte2 "mond/wind/utils/byte"
	"reflect"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	bytePool = byte2.NewBufferPool(byte2.BufferPoolOption{InitBufferSize: 1024})
	mr       = jsonpb.Marshaler{EmitDefaults: false, Indent: "", OrigName: true, EnumsAsInts: true}
)

func GrpcServerMiddleware() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log := GetLogger()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		defer func() {
			buffer := bytePool.Get()
			defer bytePool.Put(buffer)
			useTime := utils.UseTimeToStr(time.Now().Sub(start))
			//useTime := fmt.Sprintf("%.2f", float64(time.Now().UnixNano()-start.UnixNano())/float64(time.Microsecond))
			fields := []zap.Field{}
			if req != nil {
				mr.Marshal(buffer, req.(proto.Message))
				reqMap := map[string]interface{}{}
				json.Unmarshal(buffer.Bytes(), &reqMap)
				fields = append(fields, zap.Any("req", reqMap))
			}
			//fmt.Println(resp != nil)
			//fmt.Println(reflect.ValueOf(resp).IsNil())
			if err == nil && !reflect.ValueOf(resp).IsNil() {
				buffer.Reset()
				mr.Marshal(buffer, resp.(proto.Message))
				respMap := map[string]interface{}{}
				json.Unmarshal(buffer.Bytes(), &respMap)
				fields = append(fields, zap.Any("resp", respMap))
			} else {
				fields = append(fields, zap.Any("resp", nil))
			}
			fields = append(fields, zap.String("type", "grpcServer"),
				zap.String("method", info.FullMethod),
				zap.String("useTime", fmt.Sprintf("%s", useTime)))

			//reqBs, _ := json.Marshal(req)
			//respBs, _ := json.Marshal(resp)
			//reqMap := map[string]interface{}{}
			//respMap := map[string]interface{}{}
			//json.Unmarshal(reqBs, &reqMap)
			//json.Unmarshal(respBs, &respMap)

			if err != nil {
				fields = append(fields, zap.Any("err", err))
				log.Error(ctx, "", fields...)
			} else {
				log.Info(ctx, "", fields...)
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func GrpcClientMiddleware(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	var err error

	defer func() {
		fields := []zap.Field{}
		buffer := bytePool.Get()
		defer bytePool.Put(buffer)

		if req != nil {
			buffer.Reset()
			mr.Marshal(buffer, req.(proto.Message))
			reqMap := map[string]interface{}{}
			json.Unmarshal(buffer.Bytes(), &reqMap)
			fields = append(fields, zap.Any("req", reqMap))
		}

		if err == nil && !reflect.ValueOf(reply).IsNil() {
			buffer.Reset()
			mr.Marshal(buffer, reply.(proto.Message))
			respMap := map[string]interface{}{}
			json.Unmarshal(buffer.Bytes(), &respMap)
			fields = append(fields, zap.Any("resp", respMap))
		} else {
			fields = append(fields, zap.Any("resp", nil))
		}

		//reqBs, _ := json.Marshal(req)
		//respBs, _ := json.Marshal(reply)
		//reqMap := map[string]interface{}{}
		//respMap := map[string]interface{}{}
		//json.Unmarshal(reqBs, &reqMap)
		//json.Unmarshal(respBs, &respMap)
		useTime := utils.UseTimeToStr(time.Now().Sub(start))
		//fields := []zap.Field{
		//	zap.Any("req", reqMap),
		//	zap.Any("resp", respMap),
		//	zap.String("type", "grpcClient"),
		//	zap.String("method", method),
		//	//zap.String("traceId", mctx.GetTraceId(ctx)),
		//	zap.String("targetHostName", utils.ExtractCtxTargetHostName(ctx)),
		//	zap.String("useTime", useTime),
		//}
		fields = append(fields, []zap.Field{
			//zap.Any("req", reqMap),
			//zap.Any("resp", respMap),
			zap.String("type", "grpcClient"),
			zap.String("method", method),
			zap.String("targetHostName", utils.ExtractCtxTargetHostName(ctx)),
			zap.String("useTime", useTime),
		}...)
		if err != nil {
			e := merr.ParseErrToMetaErr(err)
			fields = append(fields, zap.Any("err", e))
			GetLogger().Error(ctx, "发起的请求", fields...)
		} else {
			GetLogger().Info(ctx, "发起的请求", fields...)
		}
	}()
	err = invoker(ctx, method, req, reply, cc, opts...)
	return err
}
