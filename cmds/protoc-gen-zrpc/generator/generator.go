package generator

import (
	"google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

func genFile(gen *protogen.Plugin, f *protogen.File) {
	genGo(gen, f)
	genRpc(gen, f)
}

func genGo(gen *protogen.Plugin, f *protogen.File) {
	internal_gengo.GenerateFile(gen, f)
}

func genRpc(gen *protogen.Plugin, f *protogen.File) {
}
