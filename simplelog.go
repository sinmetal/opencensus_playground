package main

import (
	"context"
	"fmt"
	"time"

	"go.opencensus.io/trace"
)

type SimpleLogService struct {
	BQ *BigQueryService
}

func NewSimpleLogService(bq *BigQueryService) *SimpleLogService {
	return &SimpleLogService{
		BQ: bq,
	}
}

func (s *SimpleLogService) Output(ctx context.Context, body string, goroutineNumber int) (rerr error) {
	const name = "fmt.Printf"
	ctx, span := StartSpan(ctx, name)
	defer func() {
		if rerr != nil {
			span.SetStatus(trace.Status{trace.StatusCodeInternal, rerr.Error()})
		}
		span.End()
	}()
	logSize := int64(len(body))
	span.AddAttributes(trace.Int64Attribute("bodySize", logSize))
	if err := RecordMeasurement("fmt.Printf", logSize); err != nil {
		fmt.Printf("failed RecordMeasurement. err=%+v\n", err)
	}

	defer func(n time.Time) {
		d := time.Since(n)
		if d.Seconds() > 1 {
			fmt.Printf("go:%d:fmt.Printf:WriteLogTime:%v/ bodySize=%v \n", goroutineNumber, d, int64(len(body)))
		}
		if s.BQ != nil {
			if err := s.BQ.Insert(ctx, &BQRow{Name: name, LogSize: logSize, Latency: d.Seconds()}); err != nil {
				panic(err)
			}
		}
	}(time.Now())

	fmt.Printf("fmt.Printf:%s\n", body)

	return nil
}
