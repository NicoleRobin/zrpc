package generator

import (
	"github.com/nicolerobin/zrpc/cmds/protoc-gen-zrpc/generator/cg"
	"google.golang.org/protobuf/compiler/protogen"
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

func (s *serviceRender) render() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderClient() cg.Builder {
	return cg.ComposeBuilder{}
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
