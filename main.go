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

	zapLogger, err := NewZapLogger()
	if err != nil {
		panic(err)
	}
	zapTommy351Logger, err := NewTommy351ZapLog()
	if err != nil {
		panic(err)
	}

	var text string
	for {
		sample(context.Background())

		if err := zapLogger.Write(context.Background(), text); err != nil {
			fmt.Printf("failed zapLogger.Write len=%+v,err=%+v\n", len(text), err)
			text = ""
		}

		if err := zapTommy351Logger.Write(context.Background(), text); err != nil {
			fmt.Printf("failed zapTommy351Logger.Write len=%+v,err=%+v\n", len(text), err)
			text = ""
		}

		time.Sleep(1*time.Second + time.Duration(rand.Intn(10))*time.Second)
		text += RandString(1024)
	}
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
