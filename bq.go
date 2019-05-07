package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

type BigQueryService struct {
	BQ *bigquery.Client
}

func NewBigQueryService(ctx context.Context, projectID string) (*BigQueryService, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &BigQueryService{
		BQ: client,
	}, nil
}

type BQRow struct {
	Name      string
	LogSize   int64
	Latency   float64
	Timestamp time.Time
}

func (s *BigQueryService) Insert(ctx context.Context, row *BQRow) error {
	inserter := s.BQ.Dataset("pottary").Table("OpenCensus").Inserter()
	rows := []*BQRow{
		row,
	}
	if err := inserter.Put(ctx, rows); err != nil {
		switch e := err.(type) {
		case bigquery.PutMultiError:
			rowInsertionError := e[0] // the first failed row
			for _, err := range rowInsertionError.Errors {
				log.Printf("err = %v", err)
			}
		}
		return err
	}
	return nil
}
