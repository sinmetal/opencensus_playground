package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

type BigQueryService struct {
	BQ         *bigquery.Client
	bufferChan chan *BQRow
}

func NewBigQueryService(ctx context.Context, projectID string) (*BigQueryService, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	bq := &BigQueryService{
		BQ:         client,
		bufferChan: make(chan *BQRow, 10000),
	}

	return bq, nil
}

type BQRow struct {
	Name      string
	LogSize   int64
	Latency   float64
	Timestamp time.Time
}

func (s *BigQueryService) Insert(ctx context.Context, row *BQRow) error {
	s.bufferChan <- row
	return nil
}

func (s *BigQueryService) InfinityFlush() {
	inserter := s.BQ.Dataset("pottary").Table("OpenCensus").Inserter()
	rows := []*BQRow{}
	count := 0
	for r := range s.bufferChan {
		rows = append(rows, r)
		count++
		if count > 1000 {
			if err := inserter.Put(context.Background(), rows); err != nil {
				switch e := err.(type) {
				case bigquery.PutMultiError:
					rowInsertionError := e[0] // the first failed row
					for _, err := range rowInsertionError.Errors {
						log.Printf("err = %v", err)
					}
				}
			}
			rows = []*BQRow{}
			count = 0
		}
	}

}
