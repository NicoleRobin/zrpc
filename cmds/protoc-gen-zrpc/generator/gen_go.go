package generator

import (
	"context"

	"github.com/nicolerobin/zrpc/log"
	"google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

func genGo(ctx context.Context, gen *protogen.Plugin, file *protogen.File) {
	log.Info(ctx, "entrance")
	internal_gengo.GenerateFile(gen, file)
}
