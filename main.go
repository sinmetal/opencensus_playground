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

	for {
		sample(context.Background())
		time.Sleep(3 * time.Second)
	}
}

func sample(ctx context.Context) {
	ctx, span := StartSpan(ctx, "sample")
	defer span.End()

	r := rand.Intn(10)
	if r > 5 {
		internalSample(ctx)
	}

	fmt.Println(time.Now())
}

func internalSample(ctx context.Context) {
	ctx, span := StartSpan(ctx, "internalSample", trace.WithSampler(trace.AlwaysSample()))
	defer span.End()

	fmt.Println(time.Now())
}
