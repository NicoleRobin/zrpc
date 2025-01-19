package main

import (
	"context"

	"github.com/nicolerobin/zrpc/log"
)

func main() {
	ctx := context.Background()
	log.SetLevel(log.LevelDebug)
	log.Debug(ctx, "debug")
	log.Info(ctx, "info")
	log.Warn(ctx, "warn")
	log.Error(ctx, "error")
}
