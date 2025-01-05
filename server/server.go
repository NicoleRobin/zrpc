package server

import "google.golang.org/grpc"

type rpcServiceRegister interface {
	RegisterService(s grpc.ServiceRegistrar, handler interface{})
}

func RegisterService(r rpcServiceRegister, handler interface{}) {
}
