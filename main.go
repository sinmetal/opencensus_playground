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
	if gcpmetadata.OnGCP() {
		project, err := gcpmetadata.GetProjectID()
		if err != nil {
			panic(err)
		}

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
	}

	zapLogger, err := NewZapLogger()
	if err != nil {
		panic(err)
	}
	zapTommy351Logger, err := NewTommy351ZapLog()
	if err != nil {
		panic(err)
	}

	endCh := make(chan error, 10)
	for i := 0; i < 1000; i++ {
		go func() {
			var text string
			var texts []string
			for {
				// sample(context.Background())

				if err := zapLogger.Write(context.Background(), text); err != nil {
					fmt.Printf("failed zapLogger.Write len=%+v,err=%+v\n", len(text), err)
				}

				if err := zapLogger.WriteNewLine(context.Background(), texts); err != nil {
					fmt.Printf("failed zapLogger.WriteNewLine len=%+v,err=%+v\n", len(texts), err)
				}

				if err := zapTommy351Logger.Write(context.Background(), text); err != nil {
					fmt.Printf("failed zapTommy351Logger.Write len=%+v,err=%+v\n", len(text), err)
				}

				if err := zapTommy351Logger.WriteNewLine(context.Background(), texts); err != nil {
					fmt.Printf("failed zapTommy351Logger.WriteNewLine len=%+v,err=%+v\n", len(texts), err)
				}

				text += RandString(1024)
				for i := 0; i < 10; i++ {
					texts = append(texts, RandString(128))
				}
			}
		}()
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
