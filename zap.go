package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func WriteZapLog(ctx context.Context, body string) (rerr error) {
	ctx, span := StartSpan(ctx, "writeZapLog")
	defer span.End()

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	defer func() {
		err := logger.Sync() // flushes buffer, if any
		if rerr == nil {
			rerr = err
			return
		}
		fmt.Printf("failed WriteZapLog Sync. err=%+v\n", err)
	}()
	sugar := logger.Sugar()
	sugar.Infof("WriteZapLog %s", body)

	return nil
}
