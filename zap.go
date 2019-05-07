package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tommy351/zap-stackdriver"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
	bq     *BigQueryService
}

func NewZapLogger(bq *BigQueryService) (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{
		logger: logger,
		bq:     bq,
	}, nil
}

func (l *ZapLogger) Write(ctx context.Context, body string, goroutineNumber int) (rerr error) {
	const name = "zap.Write"
	ctx, span := StartSpan(ctx, name)
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	logSize := int64(len(body))
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))
	if err := RecordMeasurement("zap.Write", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:zap.Write:WriteZapTime:%v/ bodySize=%v \n", goroutineNumber, d, int64(len(body)))
		}
		if l.bq != nil {
			if err := l.bq.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds(), Timestamp: time.Now()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())
	// localでは動くけど、GKE上では `sync /dev/stderr: invalid argument` と言われてエラーになる
	//defer func() {
	//	err := l.logger.Sync() // flushes buffer, if any
	//	if rerr == nil {
	//		rerr = err
	//		return
	//	}
	//	fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	//}()
	sugar := l.logger.Sugar()
	sugar.Infof("WriteZapLog:%s", body)

	return nil
}

func (l *ZapLogger) WriteNoSugar(ctx context.Context, body string, goroutineNumber int) (rerr error) {
	const name = "zap.WriteNoSugar"
	ctx, span := StartSpan(ctx, name)
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	logSize := int64(len(body))
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))
	if err := RecordMeasurement("zap.WriteNoSugar", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:zap.WriteNoSugar:WriteZapTime:%v/ bodySize=%v \n", goroutineNumber, d, int64(len(body)))
		}
		if l.bq != nil {
			if err := l.bq.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds(), Timestamp: time.Now()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())
	// localでは動くけど、GKE上では `sync /dev/stderr: invalid argument` と言われてエラーになる
	//defer func() {
	//	err := l.logger.Sync() // flushes buffer, if any
	//	if rerr == nil {
	//		rerr = err
	//		return
	//	}
	//	fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	//}()
	l.logger.Info(fmt.Sprintf("WriteZapLog:%s", body))

	return nil
}

func (l *ZapLogger) WriteNewLine(ctx context.Context, body []string, goroutineNumber int) (rerr error) {
	const name = "zap.WriteNewLine"
	ctx, span := StartSpan(ctx, name)
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	span.AddAttributes(trace.Int64Attribute("bodyLength", int64(len(body))))

	text := strings.Join(body, "\n")
	logSize := int64(len(text))
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))
	if err := RecordMeasurement("zap.WriteNewLine", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:zap.WriteNewLine:WriteZapTime:%v/ bodySize=%v \n", goroutineNumber, d, logSize)
		}
		if l.bq != nil {
			if err := l.bq.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds(), Timestamp: time.Now()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())

	// localでは動くけど、GKE上では `sync /dev/stderr: invalid argument` と言われてエラーになる
	//defer func() {
	//	err := l.logger.Sync() // flushes buffer, if any
	//	if rerr == nil {
	//		rerr = err
	//		return
	//	}
	//	fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	//}()
	sugar := l.logger.Sugar()
	sugar.Infof("WriteZapLog:%s", text)

	return nil
}

func (l *ZapLogger) WriteNewLineNoSugar(ctx context.Context, body []string, goroutineNumber int) (rerr error) {
	const name = "zap.WriteNewLineNoSugar"
	ctx, span := StartSpan(ctx, name)
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	span.AddAttributes(trace.Int64Attribute("bodyLength", int64(len(body))))

	text := strings.Join(body, "\n")
	logSize := int64(len(text))
	if err := RecordMeasurement("zap.WriteNewLineNoSugar", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:zap.WriteNewLineNoSugar:WriteZapTime:%v/ bodySize=%v \n", goroutineNumber, d, int64(len(body)))
		}
		if l.bq != nil {
			if err := l.bq.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds(), Timestamp: time.Now()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())

	// localでは動くけど、GKE上では `sync /dev/stderr: invalid argument` と言われてエラーになる
	//defer func() {
	//	err := l.logger.Sync() // flushes buffer, if any
	//	if rerr == nil {
	//		rerr = err
	//		return
	//	}
	//	fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	//}()
	l.logger.Info(fmt.Sprintf("WriteZapLog:%s", text))

	return nil
}

type Tommy351ZapLogger struct {
	logger *zap.Logger
	bq     *BigQueryService
}

func NewTommy351ZapLog(bq *BigQueryService) (*Tommy351ZapLogger, error) {
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
		bq:     bq,
	}, nil
}

func (l *Tommy351ZapLogger) Write(ctx context.Context, body string, goroutineNumber int) (rerr error) {
	const name = "tommy351ZapLogger.Write"
	ctx, span := StartSpan(ctx, "tommy351ZapLogger.Write")
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	logSize := int64(len(body))
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))
	if err := RecordMeasurement("tommy351ZapLogger.Write", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:tommy351ZapLogger.Write:WriteZapTime:%v/ bodySize=%v \n", goroutineNumber, d, int64(len(body)))
		}
		if l.bq != nil {
			if err := l.bq.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds(), Timestamp: time.Now()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())

	// localでは動くけど、GKE上では `sync /dev/stderr: invalid argument` と言われてエラーになる
	//defer func() {
	//	err := l.logger.Sync() // flushes buffer, if any
	//	if rerr == nil {
	//		rerr = err
	//		return
	//	}
	//	fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	//}()
	sugar := l.logger.Sugar()
	sugar.Infof("WriteZapLog:%s", body)

	return nil
}

func (l *Tommy351ZapLogger) WriteNoSugar(ctx context.Context, body string, goroutineNumber int) (rerr error) {
	const name = "tommy351ZapLogger.WriteNoSugar"
	ctx, span := StartSpan(ctx, name)
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	logSize := int64(len(body))
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))
	if err := RecordMeasurement("tommy351ZapLogger.WriteNoSugar", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:tommy351ZapLogger.WriteNoSugar:WriteZapTime:%v/ bodySize=%v \n", goroutineNumber, d, int64(len(body)))
		}
		if l.bq != nil {
			if err := l.bq.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds(), Timestamp: time.Now()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())

	// localでは動くけど、GKE上では `sync /dev/stderr: invalid argument` と言われてエラーになる
	//defer func() {
	//	err := l.logger.Sync() // flushes buffer, if any
	//	if rerr == nil {
	//		rerr = err
	//		return
	//	}
	//	fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	//}()
	l.logger.Info(fmt.Sprintf("WriteZapLog:%s", body))

	return nil
}

func (l *Tommy351ZapLogger) WriteNewLine(ctx context.Context, body []string, goroutineNumber int) (rerr error) {
	const name = "tommy351ZapLogger.WriteNewLine"
	ctx, span := StartSpan(ctx, name)
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	span.AddAttributes(trace.Int64Attribute("bodyLength", int64(len(body))))

	text := strings.Join(body, "\n")
	logSize := int64(len(text))
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))
	if err := RecordMeasurement("tommy351ZapLogger.WriteNewLine", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:tommy351ZapLogger.WriteNewLine:WriteZapTime:%v/ bodySize=%v \n", goroutineNumber, d, int64(len(body)))
		}
		if l.bq != nil {
			if err := l.bq.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds(), Timestamp: time.Now()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())

	// localでは動くけど、GKE上では `sync /dev/stderr: invalid argument` と言われてエラーになる
	//defer func() {
	//	err := l.logger.Sync() // flushes buffer, if any
	//	if rerr == nil {
	//		rerr = err
	//		return
	//	}
	//	fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	//}()
	sugar := l.logger.Sugar()
	sugar.Infof("WriteZapLog:%s", text)

	return nil
}

func (l *Tommy351ZapLogger) WriteNewLineNoSugar(ctx context.Context, body []string, goroutineNumber int) (rerr error) {
	const name = "tommy351ZapLogger.WriteNewLineNoSugar"
	ctx, span := StartSpan(ctx, name)
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	span.AddAttributes(trace.Int64Attribute("bodyLength", int64(len(body))))

	text := strings.Join(body, "\n")
	logSize := int64(len(text))
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))
	if err := RecordMeasurement("tommy351ZapLogger.WriteNewLineNoSugar", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:tommy351ZapLogger.WriteNewLineNoSugar:WriteZapTime:%v/ bodySize=%v \n", goroutineNumber, d, int64(len(body)))
		}
		if l.bq != nil {
			if err := l.bq.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds(), Timestamp: time.Now()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())

	// localでは動くけど、GKE上では `sync /dev/stderr: invalid argument` と言われてエラーになる
	//defer func() {
	//	err := l.logger.Sync() // flushes buffer, if any
	//	if rerr == nil {
	//		rerr = err
	//		return
	//	}
	//	fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	//}()
	l.logger.Info(fmt.Sprintf("WriteZapLog:%s", body))

	return nil
}
