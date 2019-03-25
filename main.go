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
		// trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	var text string
	for {
		sample(context.Background())
		if err := WriteZapLog(context.Background(), text); err != nil {
			fmt.Printf("failed WriteZapLog len=%+v,err=%+v\n", len(text), err)
			text = ""
		}
		time.Sleep(3 * time.Second)
		text += "a"
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
	ctx, span := StartSpan(ctx, "internalSample", trace.WithSampler(trace.AlwaysSample()))
	defer span.End()

	time.Sleep(100 * time.Millisecond)
}
