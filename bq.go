package main

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/google/uuid"
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
	Name    string
	LogSize int64
	Latency float64
}

func (r *BQRow) Save() (row map[string]bigquery.Value, insertID string, err error) {
	m := map[string]bigquery.Value{}
	m["Name"] = r.Name
	m["LogSize"] = r.LogSize
	m["Latency"] = r.Latency
	return m, uuid.New().String(), nil
}

func (s *BigQueryService) Insert(ctx context.Context, row *BQRow) error {
	inserter := s.BQ.Dataset("pottary").Table("Latency").Inserter()
	if err := inserter.Put(ctx, row); err != nil {
		return err
	}
	return nil
}
