package generator

type rpcMethodInfo struct {
	name                 string
	protoMethodName      string
	methodName           string
	methodPath           string
	handlerName          string
	reqType              string
	resType              string
	streamIndex          int
	streamParamName      string
	serverStreamName     string
	serverStreamTypeName string
	isClientStream       bool
	isServerStream       bool
}
