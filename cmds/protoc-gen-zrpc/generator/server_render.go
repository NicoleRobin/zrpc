package generator

import (
	"github.com/nicolerobin/zrpc/cmds/protoc-gen-zrpc/generator/cg"
	"google.golang.org/protobuf/compiler/protogen"
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
	return cg.ComposeBuilder{}
}
func (s *serviceRender) renderTypeInfo() cg.Builder {
	items := make([]cg.Builder, len(s.rpcInfo))
	for _, rpcInfo := range s.rpcInfo {
		items = append(items, cg.RawString(rpcInfo.name))
	}

	return cg.Const(items)
}
