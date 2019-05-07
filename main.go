package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/sinmetal/gcpmetadata"
	"go.opencensus.io/trace"
)

func main() {
	var bq *BigQueryService
	if gcpmetadata.OnGCP() {
		project, err := gcpmetadata.GetProjectID()
		if err != nil {
			panic(err)
		}
		fmt.Printf("GOOGLE_CLOUD_PROJECT:%s\n", project)

		{
			exporter, err := stackdriver.NewExporter(stackdriver.Options{
				ProjectID: project,
			})
			if err != nil {
				panic(err)
			}
			trace.RegisterExporter(exporter)
			trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		}
		{
			exporter := InitExporter()
			InitOpenCensusStats(exporter)
		}

		// BQに出力する処理に時間がかかって、ログ出力量が減り、つまらなくなるので止めた
		//bq, err = NewBigQueryService(context.Background(), project)
		//if err != nil {
		//	panic(err)
		//}
		//
		//go func() {
		//	bq.InfinityFlush()
		//}()
	}

	slogger := NewSimpleLogService(bq)

	zapLogger, err := NewZapLogger(bq)
	if err != nil {
		panic(err)
	}
	zapTommy351Logger, err := NewTommy351ZapLog(bq)
	if err != nil {
		panic(err)
	}

	endCh := make(chan error, 10)
	for i := 0; i < 1000; i++ {
		i := i
		go func(i int) {
			var text string
			var texts []string
			for {
				// sample(context.Background())

				if err := slogger.Output(context.Background(), text, i); err != nil {
					fmt.Printf("failed outputSimpleLog len=%+v,err=%+v\n", len(text), err)
				}

				if err := zapLogger.Write(context.Background(), text, i); err != nil {
					fmt.Printf("failed zapLogger.Write len=%+v,err=%+v\n", len(text), err)
				}

				if err := zapLogger.WriteNoSugar(context.Background(), text, i); err != nil {
					fmt.Printf("failed zapLogger.WriteNoSugar len=%+v,err=%+v\n", len(text), err)
				}

				if err := zapLogger.WriteNewLine(context.Background(), texts, i); err != nil {
					fmt.Printf("failed zapLogger.WriteNewLine len=%+v,err=%+v\n", len(texts), err)
				}

				if err := zapLogger.WriteNewLineNoSugar(context.Background(), texts, i); err != nil {
					fmt.Printf("failed zapLogger.WriteNewLineNoSugar len=%+v,err=%+v\n", len(texts), err)
				}

				if err := zapTommy351Logger.Write(context.Background(), text, i); err != nil {
					fmt.Printf("failed zapTommy351Logger.Write len=%+v,err=%+v\n", len(text), err)
				}

				if err := zapTommy351Logger.WriteNoSugar(context.Background(), text, i); err != nil {
					fmt.Printf("failed zapTommy351Logger.WriteNoSugar len=%+v,err=%+v\n", len(text), err)
				}

				if err := zapTommy351Logger.WriteNewLine(context.Background(), texts, i); err != nil {
					fmt.Printf("failed zapTommy351Logger.WriteNewLine len=%+v,err=%+v\n", len(texts), err)
				}

				if err := zapTommy351Logger.WriteNewLineNoSugar(context.Background(), texts, i); err != nil {
					fmt.Printf("failed zapTommy351Logger.WriteNewLineNoSugar len=%+v,err=%+v\n", len(texts), err)
				}

				text += RandString(1024)
				for j := 0; j < 10; j++ {
					texts = append(texts, RandString(128))
				}
			}
		}(i)
	}

	err = <-endCh
	fmt.Printf("BOMB %+v", err)
}

func sample(ctx context.Context) {
	ctx, span := StartSpan(ctx, "sample")
	defer span.End()

	r := rand.Intn(10)
	if r > 5 {
		internalSample(ctx)
	}
}

func internalSample(ctx context.Context) {
	ctx, span := StartSpan(ctx, "internalSample")
	defer span.End()

	time.Sleep(30 * time.Millisecond)
}
