package logger

import (
	"context"
	"mond/wind/env"
	mctx "mond/wind/utils/ctx"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
}

type logger struct {
	_logger *zap.Logger
}

func getContextField(ctx context.Context) []zap.Field {
	return []zap.Field{
		zap.String("appId", env.GetAppId()),
		zap.String("traceId", mctx.GetTraceId(ctx)),
	}
}
func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getContextField(ctx)...)
	l._logger.Debug(msg,
		fields...,
	)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getContextField(ctx)...)
	l._logger.Info(msg,
		fields...,
	)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, getContextField(ctx)...)
	l._logger.Error(msg,
		fields...,
	)
}

var (
	once    sync.Once
	_logger *logger
)

func GetLogger() Logger {
	if _logger == nil {
		once.Do(func() {
			conf := zap.NewProductionConfig()
			conf.Sampling = nil
			//FIXME:  为了调试
			conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			conf.Encoding = "json"
			conf.EncoderConfig.TimeKey = "time"
			conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			l, _ := conf.Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.PanicLevel))
			_logger = &logger{
				_logger: l,
			}
		})
	}
	return _logger
}
