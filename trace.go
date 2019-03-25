package main

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"
)

func StartSpan(ctx context.Context, name string, o ...trace.StartOption) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, fmt.Sprintf("/oc/%s", name), o...)
}
