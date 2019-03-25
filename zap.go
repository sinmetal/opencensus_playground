package main

import (
	"context"
	"fmt"

	"github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{
		logger: logger,
	}, nil
}

func (l *ZapLogger) Write(ctx context.Context, body string) (rerr error) {
	ctx, span := StartSpan(ctx, "zap.Write")
	defer span.End()

	defer func() {
		err := l.logger.Sync() // flushes buffer, if any
		if rerr == nil {
			rerr = err
			return
		}
		fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	}()
	sugar := l.logger.Sugar()
	sugar.Infof("WriteZapLog:%s", body)

	return nil
}

type Tommy351ZapLogger struct {
	logger *zap.Logger
}

func NewTommy351ZapLog() (*Tommy351ZapLogger, error) {
	config := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:         "json",
		EncoderConfig:    stackdriver.EncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &stackdriver.Core{
			Core: core,
		}
	}), zap.Fields(
		stackdriver.LogServiceContext(&stackdriver.ServiceContext{
			Service: "foo",
			Version: "bar",
		}),
	))

	if err != nil {
		return nil, err
	}

	return &Tommy351ZapLogger{
		logger: logger,
	}, nil
}

func (l *Tommy351ZapLogger) Write(ctx context.Context, body string) (rerr error) {
	ctx, span := StartSpan(ctx, "tommy351ZapLogger.Write")
	defer span.End()

	defer func() {
		err := l.logger.Sync() // flushes buffer, if any
		if rerr == nil {
			rerr = err
			return
		}
		fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	}()

	return nil
}
