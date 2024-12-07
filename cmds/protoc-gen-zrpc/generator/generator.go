package generator

import (
	"github.com/nicolerobin/zrpc/cmds/protoc-gen-zrpc/constant"
	"google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func Generate() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if f.Generate {
				genFile(gen, f)
			}
		}
		return nil
	})
}

func genFile(gen *protogen.Plugin, f *protogen.File) {
	genGo(gen, f)
	genRpc(gen, f)
}

func genGo(gen *protogen.Plugin, file *protogen.File) {
	internal_gengo.GenerateFile(gen, file)
}

func genRpc(gen *protogen.Plugin, file *protogen.File) {
	if len(file.Services) == 0 {
		return
	}

	name := file.GeneratedFilenamePrefix + "_grpc.pb.go"
	g := gen.NewGeneratedFile(name, file.GoImportPath)

	g.P("// Code generated by protoc-gen-zrpc. DO NOT EDIT.")
	g.P("//")
	g.P("// version: protoc-gen-zrpc ", constant.Version)
	g.P("// source: ", file.Desc.Path())
	g.P()
	g.P("package ", file.GoPackageName)

	pkg := string(file.Desc.Package())

	for _, service := range file.Services {
		r := serviceRender{
			service:  service,
			pkgName:  pkg,
			fileName: service.Desc.ParentFile().Path(),
		}

		g.P(r.render(file))
		g.P()
	}
}

type serviceRender struct {
	service  *protogen.Service
	pkgName  string
	fileName string
}

func (s *serviceRender) render(file *protogen.File) string {
	return ""
}
