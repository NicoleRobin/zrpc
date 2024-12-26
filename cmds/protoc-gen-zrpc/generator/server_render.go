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
	name := "Unimplemented" + s.serverInterfaceName

	m := cg.Method("", cg.P(name))
	methods := make(cg.ComposeBuilder, 0, len(s.rpcInfo))

	var (
		ctx, errFunc cg.RawString
		code         string
	)

	if len(s.rpcInfo) != 0 {
		ctx = cg.S(s.qualified(contextPackage.Ident("Context")))
		errFunc = cg.S(s.qualified(grpcStatusPackage.Ident("Errorf")))
		code = s.qualified(grpcCodesPackage.Ident("Unimplemented"))
	}

	for _, rpcInfo := range s.rpcInfo {
		m = m.Name(rpcInfo.methodName)

		err := errFunc.Call(code, strconv.Quote("method"+rpcInfo.methodName+" not implemented"))

		if rpcInfo.isClientStream || rpcInfo.isServerStream {
			m = m.Return(cg.S("error")).Body(cg.Return(err))

			if !rpcInfo.isClientStream {
				m = m.Param(
					cg.Param("ctx", ctx),
					cg.Param("req", cg.P(rpcInfo.reqType)),
					cg.Param("stream", cg.S(rpcInfo.serverStreamName)),
				)
			} else {
				m = m.Param(
					cg.Param("ctx", ctx),
					cg.Param("stream", cg.S(rpcInfo.serverStreamTypeName)),
				)
			}
		} else {
			m = m.Param(cg.Param("ctx", ctx), cg.Param(
				"req", cg.P(rpcInfo.reqType),
			)).Return(cg.P(rpcInfo.reqType), cg.S("error")).Body(cg.Return(cg.S("nil"), err))
		}
		methods = append(methods, m)
	}

	return cg.ComposeBuilder{
		cg.Struct(name),
		methods,
	}
}

func (s *serviceRender) renderServiceDesc() cg.Builder {
	methods := make([]cg.Builder, 0, len(s.rpcInfo))
	streams := make([]cg.Builder, 0)

	for _, rpcInfo := range s.rpcInfo {
		b := cg.StructLiteral("").Body(cg.KV("Handler", cg.S(rpcInfo.handlerName)))
		method := strconv.Quote(rpcInfo.protoMethodName)

		if rpcInfo.isServerStream || rpcInfo.isClientStream {
			b = b.AppendBody(
				cg.KV("StreamName", cg.S(method)),
				cg.KV("ClientStream", cg.B(rpcInfo.isClientStream)),
				cg.KV("ServerStream", cg.B(rpcInfo.isServerStream)),
			)

			streams = append(streams, b)
		} else {
			b = b.AppendBody(cg.KV("MethodName", cg.S(method)))
			methods = append(methods, b)
		}
	}
	return cg.Var(s.serviceDesc).Value(
		cg.StructLiteral(s.qualified(grpcPackage.Ident("ServiceDesc"))).Body(
			cg.KV("ServiceName", cg.S(strconv.Quote(s.serviceFullName))),
			cg.KV("HandlerType", cg.Paren(cg.P(s.serverInterfaceName)).Call("nil")),
			cg.KV("Methods", cg.Array(s.qualified(grpcPackage.Ident("MethodDesc"))).Body(methods...)),
			cg.KV("Stream", cg.Array(s.qualified(grpcPackage.Ident("StreamDesc"))).Body(streams...)),
			cg.KV("Metadata", cg.S(strconv.Quote(s.fileName))),
		))
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
