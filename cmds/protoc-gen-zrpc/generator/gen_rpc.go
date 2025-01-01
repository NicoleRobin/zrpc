package generator

import (
	"context"

	"github.com/nicolerobin/zrpc/cmds/protoc-gen-zrpc/constant"
	"github.com/nicolerobin/zrpc/log"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	errorsPackage       = protogen.GoImportPath("errors")
	contextPackage      = protogen.GoImportPath("context")
	ioPackage           = protogen.GoImportPath("io")
	syncPackage         = protogen.GoImportPath("sync")
	grpcPackage         = protogen.GoImportPath("google.golang.org/grpc")
	grpcStatusPackage   = protogen.GoImportPath("google.golang.org/grpc/status")
	grpcCodesPackage    = protogen.GoImportPath("google.golang.org/grpc/codes")
	grpcMetadataPackage = protogen.GoImportPath("google.golang.org/grpc/metadata")
	protoPackage        = protogen.GoImportPath("google.golang.org/protobuf/proto")
	clientPackage       = protogen.GoImportPath("github.com/nicolerobin/zrpc/client")
	serverPackage       = protogen.GoImportPath("github.com/nicolerobin/zrpc/server")
	rpcPackage          = protogen.GoImportPath("github.com/nicolerobin/zrpc/core/rpc")
)

func genRpc(ctx context.Context, gen *protogen.Plugin, file *protogen.File) {
	log.Info(ctx, "entrance")
	if len(file.Services) == 0 {
		return
	}

	name := file.GeneratedFilenamePrefix + "_grpc.pb.go"
	generatedFile := gen.NewGeneratedFile(name, file.GoImportPath)

	generatedFile.P("// Code generated by protoc-gen-zrpc. DO NOT EDIT.")
	generatedFile.P("//")
	generatedFile.P("// version: protoc-gen-zrpc ", constant.Version)
	generatedFile.P("// source: ", file.Desc.Path())
	generatedFile.P()
	generatedFile.P("package ", file.GoPackageName)

	pkgName := string(file.Desc.Package())
	clientGetter := generatedFile.QualifiedGoIdent(rpcPackage.Ident("ClientGetter"))
	serverRegister := generatedFile.QualifiedGoIdent(serverPackage.Ident("RegisterService"))

	for _, service := range file.Services {
		log.Infof(ctx, "loop services, service:%+v", service)
		r := serviceRender{
			service:        service,
			pkgName:        pkgName,
			fileName:       service.Desc.ParentFile().Path(),
			clientGetter:   clientGetter,
			serverRegister: serverRegister,
			qualified:      generatedFile.QualifiedGoIdent,
		}

		generatedFile.P(r.render())
		generatedFile.P()
	}
}
