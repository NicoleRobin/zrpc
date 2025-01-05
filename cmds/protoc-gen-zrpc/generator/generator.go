package generator

import (
	"context"
	"go.uber.org/zap"

	"github.com/nicolerobin/zrpc/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func Generate() {
	ctx := context.Background()
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for i, f := range gen.Files {
			log.Info(ctx, "loop files", zap.Int("i", i), zap.String("f.GeneratedFilenamePrefix", f.GeneratedFilenamePrefix))
			if f.Generate {
				genFile(ctx, gen, f)
			}
		}
		return nil
	})
}

func genFile(ctx context.Context, gen *protogen.Plugin, f *protogen.File) {
	genGo(ctx, gen, f)
	genRpc(ctx, gen, f)
}
