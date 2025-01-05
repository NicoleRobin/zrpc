package cg

import (
	"context"
	"github.com/nicolerobin/zrpc/log"
	"go.uber.org/zap"
)

type ComposeBuilder []Builder

func (b ComposeBuilder) Build() string {
	log.Info(context.Background(), "entrance")
	var s string
	for i, builder := range b {
		log.Info(context.Background(), "loop compose builder", zap.Int("i", i), zap.String("builder", builder.String()))
		s += builder.Build() + "\n"
	}

	return s
}

func (b ComposeBuilder) String() string {
	log.Info(context.Background(), "entrance")
	return b.Build()
}
