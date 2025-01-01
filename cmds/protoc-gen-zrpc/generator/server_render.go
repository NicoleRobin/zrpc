package generator

import (
	"fmt"
	"strconv"
	"strings"

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

func unexport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func (s *serviceRender) render() cg.Builder {
	return cg.ComposeBuilder{
		s.renderClient(),
		s.renderClientMethods(),
		s.renderClientInterface(),
		s.renderClientGetter(),
		s.renderServiceRegister(),
		s.renderServerHandler(),
		s.renderServerInterface(),
		s.renderUnimplemented(),
		s.renderServiceDesc(),
		s.renderRPCNames(),
		s.renderTypeInfo(),
	}
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
		).Return(cg.S(s.clientInterfaceName)).Body(
			cg.Return(cg.StructPointerLiteral(s.clientTypeName).Body(
				cg.KV("cm", cg.S(s.qualified(rpcPackage.Ident("NewRawGetter"))).Call("cc")),
			)),
		),
	}
}

func (s *serviceRender) renderClientGetter() cg.Builder {
	getterName := "Get" + s.clientInterfaceName

	gen := func(name, clientName string, args ...cg.ParamBuilder) cg.Builder {
		return cg.Var(name).Value(cg.Func("").Param(args...).Return(cg.S(s.clientInterfaceName)).Body(
			cg.Return(cg.StructPointerLiteral(s.clientTypeName).Body(
				cg.KV("cm", cg.S(s.qualified(clientPackage.Ident("Get"))).Call(clientName)),
			)),
		))
	}

	named := "GetNamed" + s.clientInterfaceName

	return cg.ComposeBuilder{
		cg.Comment(getterName + " gets the default service client named " + s.pkgName),
		gen(getterName, strconv.Quote(s.pkgName)),
		cg.Comment(named + "gets the named service client specified by name."),
		gen(named, "name", cg.Param("name", cg.S("string"))),
		genClientMocker(s.serviceName, getterName, named),
	}
}

func (s *serviceRender) renderClientMethods() cg.Builder {
	b := make(cg.ComposeBuilder, 0, len(s.service.Methods))
	s.rpcInfo = make([]rpcMethodInfo, 0, len(s.service.Methods))

	var streamCount int

	for _, m := range s.service.Methods {
		info := rpcMethodInfo{
			protoMethodName: string(m.Desc.Name()),
			methodName:      m.GoName,
			methodPath:      fmt.Sprintf("/%s/%s", s.service.Desc.FullName(), m.Desc.Name()),
			handlerName:     fmt.Sprintf("%s%sHandler", unexport(s.serviceName), m.GoName),
			name:            s.serviceName + "RPC_" + m.GoName,
			reqType:         s.qualified(m.Input.GoIdent),
			resType:         s.qualified(m.Output.GoIdent),
			isClientStream:  m.Desc.IsStreamingClient(),
			isServerStream:  m.Desc.IsStreamingServer(),
		}

		method := cg.Method("w", cg.P(s.clientTypeName)).Name(m.GoName)
		cm := method.Attr("cm")

		getClient := func(nilValue bool) cg.Builder {
			b := cg.ComposeBuilder{
				cg.Var("err").Type(cg.S("error")),
				cg.Defer(cg.Func("").Body(cm.Attr("Close").Call("err")).Call()),
				cg.DefineAssign("c, err", cm.Attr("GetClient").Call("ctx")),
			}

			var ret cg.Builder
			if nilValue {
				ret = cg.S("nil, err")
			} else {
				ret = cg.S("err")
			}

			return append(b, cg.If(cg.Ne("err", "nil")).Body(cg.Return(ret)))
		}

		var f cg.Builder
		if !info.isClientStream && !info.isServerStream {
			f = method.Param(
				cg.Param("ctx", cg.S(s.qualified(contextPackage.Ident("Context")))),
				cg.Param("in", cg.P(info.resType)),
				cg.Param("opts", cg.Variadic(s.qualified(grpcPackage.Ident("CallOption")))),
			).Return(
				cg.P(info.resType),
				cg.S("error"),
			).Body(
				getClient(true),
				cg.DefineAssign("out", cg.New(cg.S(info.resType))),
				cg.Assign("err", cg.S("c").Attr("Invoke").Call("ctx", strconv.Quote(info.methodPath), "in", "out", "opts...")),
				cg.If(cg.Ne("err", "nil")).Body(cg.Return(cg.S("nil, err"))),
				cg.Return(cg.S("out, nil")),
			)
		} else {
			info.streamIndex = streamCount
			info.streamParamName = s.serviceName + info.methodName + "Param"
			info.serverStreamName = s.serviceName + info.methodName + "Stream"
			info.serverStreamTypeName = unexport(info.serverStreamName)

			streamCount++

			f = s.renderStreamMethod(info, method, getClient(false))
		}

		s.rpcInfo = append(s.rpcInfo, info)
		b = append(b, f)
	}
	return b
}

func (s *serviceRender) renderClientInterface() cg.Builder {
	apis := make([]cg.Builder, 0, len(s.rpcInfo))

	var ctx, callOption string
	if len(s.rpcInfo) != 0 {
		ctx = s.qualified(contextPackage.Ident("Context"))
		callOption = s.qualified(grpcPackage.Ident("CallOption"))
	}

	for _, rpcInfo := range s.rpcInfo {
		var api cg.Builder
		if rpcInfo.isClientStream  || rpcInfo.isServerStream {
			api = cg.InterfaceAPI(rpcInfo.methodName).Param(
				cg.Param("ctx", cg.S(ctx)),
				cg.Param("param", cg.S(rpcInfo.streamParamName)),
				cg.Param("opts", cg.Variadic(callOption)),
			).Return(cg.S("error"))
		} else {
			api = cg.InterfaceAPI(rpcInfo.methodName).Param(
				cg.Param("ctx", cg.S(ctx)),
				cg.Param("in", cg.P(rpcInfo.resType)),
				cg.Param("opts", cg.Variadic(callOption)),
			).Return(
				cg.P(rpcInfo.resType),
				cg.S("error"),
			)
		}
		apis = append(apis, api)
	}
	return cg.Interface(s.clientInterfaceName).Body(apis...)
}

func (s *serviceRender) renderServiceRegister() cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderServerHandler() cg.Builder {
	b := make(cg.ComposeBuilder, 0, len(s.rpcInfo))

	var ctx cg.RawString
	if len(s.rpcInfo) != 0 {
		ctx = cg.S(s.qualified(contextPackage.Ident("Context")))
	}

	f := func(name string) cg.FuncBuilder {
		return cg.Func(name).Param(
			cg.Param("srv", cg.S("interface{}")),
			cg.Param("ctx", ctx),
			cg.Param("dec", cg.Func("").Param(cg.Param("", cg.S("interface{}"))).Return(cg.S("error")).AsType()),
			cg.Param("interceptor", cg.S(s.qualified(grpcPackage.Ident("UnaryServerInterceptor")))),
		).Return(cg.S("interface{}"), cg.S("error"))
	}

	for _, rpcInfo := range s.rpcInfo {
		var h cg.Builder
		if rpcInfo.isClientStream  || rpcInfo.isServerStream {
			h = cg.ComposeBuilder{
				s.renderServerStreamType(rpcInfo),
				s.renderServerStreamTypeImpl(rpcInfo),
				s.renderServerStreamHandler(rpcInfo),
			}
		} else {
			serverInfo := cg.StructPointerLiteral(s.qualified(grpcPackage.Ident("UnaryServerInfo"))).Body(
				cg.KV("Server", cg.S("srv")),
				cg.KV("FullMethod", cg.S(strconv.Quote(rpcInfo.methodPath))),
			)
		}
		b = append(b, h)
	}
	return b
}

func (s *serviceRender) renderServerInterface() cg.Builder {
	apis := make([]cg.Builder, 0, len(s.rpcInfo))

	var ctx cg.RawString
	if len(s.rpcInfo) != 0 {
		ctx = cg.S(s.qualified(contextPackage.Ident("Context")))
	}

	for _, rpcInfo := range s.rpcInfo {
		api := cg.InterfaceAPI(rpcInfo.methodName)
		if rpcInfo.isClientStream  || rpcInfo.isServerStream {
			if !rpcInfo.isClientStream {
				api = api.Param(
					cg.Param("", ctx),
					cg.Param("", cg.P(rpcInfo.reqType)),
					cg.Param("", cg.S(rpcInfo.streamParamName)),
				)
			} else {
				api = api.Param(
					cg.Param("", ctx),
					cg.Param("", cg.S(rpcInfo.serverStreamName)),
				)
			}
			api = api.Return(cg.S("error"))
		} else {
			api = api.Param(
				cg.Param("", ctx),
				cg.Param("", cg.P(rpcInfo.reqType)),
			).Return(
				cg.P(rpcInfo.resType),
				cg.S("error"),
			)
		}
		apis = append(apis, api)
	}
	return cg.Interface(s.serverInterfaceName).Body(apis...)
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

func (s *serviceRender) renderStreamMethod(info rpcMethodInfo, method cg.FuncBuilder, clientGetter cg.Builder) cg.Builder {
	errRet := cg.If(cg.S("err != nil")).Body(cg.Return(cg.S("err")))
	eof := s.qualified(ioPackage.Ident("EOF"))

	ctxType := cg.S(s.qualified(contextPackage.Ident("Context")))

	f := method.Param(
		cg.Param("ctx", ctxType),
		cg.Param("param", cg.S(info.streamParamName)),
		cg.Param("opts", cg.Variadic(s.qualified(grpcPackage.Ident("CallOption")))),
	).Return(
		cg.S("error"),
	).Body(
		clientGetter,
	)

	stream := cg.S("stream")
	param := cg.S("param")

	readLoop := cg.ComposeBuilder{
		cg.Var("mctx").Type(ctxType),
		cg.For("").Body(
			cg.DefineAssign("msg", cg.New(cg.S(info.reqType))),
			cg.Assign("err", stream.Attr("RecvMsg").Call("msg")),
			cg.If(cg.Ne("err", eof)).Body(
				cg.If(cg.Eq("err", eof)).Body(cg.Assign("err", cg.S("nil"))),
				cg.Break(),
			),
			s.renderStreamCallOnMessage(true),
		),
	}

	var body cg.Builder
	switch {
	case info.isServerStream && !info.isClientStream:
		body = cg.ComposeBuilder{
			cg.Assign("err", stream.Attr("SendMsg").Call(param.Attr("Req").String())),
			errRet,
			cg.Assign("err", stream.Attr("CloseSend").Call()),
			errRet,
			readLoop,
			cg.Return(cg.S("err")),
		}
	case !info.isServerStream && info.isClientStream:
		body = cg.ComposeBuilder{
			s.renderStreamWriteLoop(),
			errRet,
			cg.DefineAssign("msg", cg.New(cg.S(info.reqType))),
			cg.Assign("err", stream.Attr("RecvMsg").Call("msg")),
			errRet,
			cg.Var("mctx").Type(ctxType),
			s.renderStreamCallOnMessage(false),
			cg.Return(cg.S("nil")),
		}
	default:
		f = f.AppendBody(
			cg.DefineAssign("ctx, cancel", cg.S(s.qualified(contextPackage.Ident("WithCancel"))).Call("ctx")),
			cg.Defer(cg.S("cancel").Call()),
		)

		body = cg.ComposeBuilder{
			cg.DefineAssign("sendErr", cg.Make(cg.Chan(cg.S("error")), cg.S("1"))),
			cg.DefineAssign("wait", cg.Make(cg.Chan(cg.S("struct{}")))),
		}

		body = cg.ComposeBuilder{}
	}

	return cg.ComposeBuilder{
		s.renderStreamMethodParam(info),
		f.AppendBody(s.wrapClientStreamHook(info, body)),
	}
}

func (s *serviceRender) renderStreamCallOnMessage(isServerStream bool) cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) renderStreamMethodParam(info rpcMethodInfo) cg.Builder {
	return cg.ComposeBuilder{}
}

func (s *serviceRender) wrapClientStreamHook(info rpcMethodInfo, body cg.Builder) cg.Builder {
	hook := cg.S(s.qualified(grpcPackage.Ident("StreamClientHook")))
	method := strconv.Quote(info.methodName)

	return cg.ComposeBuilder{
		cg.DefineAssign("desc", cg.Attr(s.serviceDesc).Attr(cg.Index("Streams", info.streamIndex).String()),
		cg.Return(hook.Call("ctx", "c", "desc", method, cg.Func("").Param(
			cg.Param("ctx", cg.S(s.qualified(contextPackage.Ident("Context")))),
			cg.Param("stream", cg.S(s.qualified(grpcPackage.Ident("ClientStream")))),
		).Return(cg.S("error")).Body(body).String(), "opts...")),
	}
}

func (s *serviceRender) renderStreamWriteLoop() cg.Builder {
	return cg.ComposeBuilder{}
}

func genClientMocker(serviceName, getterName, namedGetterName string) cg.Builder {
	iface := serviceName + "Client"
	structName := unexport(serviceName) + "Mocker"

	method := cg.Method("m", cg.P(structName))

	genMocker := func(name, attr, getter string, param cg.ParamBuilder) cg.Builder {
		f := cg.Func("").Param(param).Return(cg.S(iface))

		return method.Name(name).Param(
			cg.Param("mocker", f.AsType()),
		).Body(
			cg.Assign(method.Attr(attr), cg.S(getter)),
			cg.Assign(cg.S(getter), f.Body(
				cg.Return(cg.Assert(cg.S("mocker").Call(param.Name()), cg.S(iface))),
			)),
		)
	}

	resetFuncBody := func(x, y string) cg.Builder {
		return cg.If(cg.Ne(x, "nil")).Body(cg.Assign(cg.S(y), cg.S(x)))
	}

	exportName := serviceName + "Mocker"

	return cg.ComposeBuilder{
		cg.Struct(structName).Body(
			cg.Param("defaultClient", cg.Func("").Return(cg.S(iface)).AsType()),
			cg.Param("namedClient", cg.Func("").Param(cg.Param("", cg.S("string"))).Return(cg.S(iface)).AsType()),
		),
		genMocker("Mock", "defaultClient", getterName, cg.NoParam),
		genMocker("MockNamed", "namedClient", namedGetterName, cg.Param("name", cg.S("string"))),
		method.Name("Reset").Body(resetFuncBody("m.defaultClient", getterName)),
		method.Name("ResetNamed").Body(resetFuncBody("m.namedClient", namedGetterName)),
		cg.Comment(exportName + " is a mocker which is used to mock " + iface + "."),
		cg.Var(exportName).Value(cg.New(cg.S(structName))),
	}
}
