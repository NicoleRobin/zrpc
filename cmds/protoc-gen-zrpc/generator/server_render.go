package generator

import (
	"github.com/nicolerobin/zrpc/cmds/protoc-gen-zrpc/generator/cg"
	"google.golang.org/protobuf/compiler/protogen"
	"strconv"
	"strings"
)

type serviceRender struct {
	service        *protogen.Service
	pkgName        string
	fileName       string
	clientGetter   string
	serverRegister string

	qualified func(protogen.GoIdent) string

	serviceFullName     string
	serviceName         string
	serviceDesc         string
	clientTypeName      string
	clientInterfaceName string
	serverInterfaceName string
	rpcInfo             []rpcMethodInfo
}

func unexport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func (s *serviceRender) render() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderClient() cg.Builder {
	s.serviceName = s.service.GoName
	s.serviceFullName = string(s.service.Desc.FullName())
	s.serverInterfaceName = s.serviceName + "Server"
	s.serviceDesc = unexport(s.serviceName) + "ServiceDesc"
	s.clientInterfaceName = s.serviceName + "Client"
	s.clientTypeName = "wrapped" + s.clientInterfaceName

	return cg.ComposeBuilder{
		cg.Struct(s.clientTypeName).Body(cg.Param("cm", cg.S(s.clientGetter))),
		cg.Func("New" + s.clientInterfaceName).Param(
			cg.Param("cc", cg.S(s.qualified(grpcPackage.Ident("ClientConnInterface")))),
		).Return(cg.S(s.clientInterfaceName)),
	}
}

func (s *serviceRender) renderClientMethods() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderClientInterface() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderClientGetter() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderServiceRegister() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderServerHandler() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderServerInterface() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderUnimplemented() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderServiceDesc() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderRPCNames() cg.Builder {
	items := make([]cg.Builder, 0, len(s.rpcInfo))

	for _, rpcInfo := range s.rpcInfo {
		items = append(items, cg.Assign(cg.S(rpcInfo.name), cg.S(strconv.Quote(rpcInfo.methodPath))))
	}
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderTypeInfo() cg.Builder {
	items := make([]cg.Builder, len(s.rpcInfo))

	var (
		fn          cg.RawString
		messageType string
		newTypeFunc cg.FuncBuilder
	)

	if len(s.rpcInfo) != 0 {
		fn = cg.S(s.qualified(rpcPackage.Ident("RegisterTypeInfo")))
		messageType = s.qualified(protoPackage.Ident("Message"))
		newTypeFunc = cg.Func("").Return(cg.S(messageType))
	}

	for _, rpcInfo := range s.rpcInfo {
		info := cg.StructLiteral(s.qualified(rpcPackage.Ident("TypeInfo"))).Body(
			cg.KV("Request", cg.Paren(cg.P(rpcInfo.reqType)).Call("nil")),
			cg.KV("Returns", cg.Paren(cg.P(rpcInfo.reqType)).Call("nil")),
			cg.KV("NewRequest", newTypeFunc.Body(cg.Return(cg.New(cg.S(rpcInfo.reqType))))),
			cg.KV("NewRequest", newTypeFunc.Body(cg.Return(cg.New(cg.S(rpcInfo.reqType))))),
		)
		items = append(items, fn.Call(rpcInfo.name, info.String()))
	}

	return cg.Func("init").Body(items...)
}
