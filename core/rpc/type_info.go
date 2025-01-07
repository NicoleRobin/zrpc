package rpc

import (
	"fmt"
	"google.golang.org/protobuf/runtime/protoiface"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type TypeInfo struct {
	Request    interface{}
	Returns    interface{}
	NewRequest interface{}
	NewReturns interface{}
}

type ParsedTypeInfo struct {
	Request    interface{}
	Returns    interface{}
	NewRequest interface{}
	NewReturns interface{}
}

var typeRegistry = make(map[string]ParsedTypeInfo)

type mixedMessage struct {
	protoiface.MessageV1
	proto.Message
}

func convertMessage(msg interface{}) proto.Message {
	switch v := msg.(type) {
	case proto.Message:
		return v
	default:
		return protoimpl.X.ProtoMessageV2Of(v)
	}
}

func convertMessageFactory(f interface{}) func() proto.Message {
	switch v := f.(type) {
	case func() protoiface.MessageV1:
		return func() proto.Message {
			v1 := v()
			v2 := protoimpl.X.ProtoMessageV2Of(v1)

			return mixedMessage{v1, v2}
		}
	default:
		return f.(func() proto.Message)
	}
}

// RegisterTypeInfo registers a type info to registry
func RegisterTypeInfo(method string, info TypeInfo) {
	if _, ok := typeRegistry[method]; ok {
		fmt.Fprintf(os.Stderr, "rpc: types already registered\n")
	}

	typeRegistry[method] = ParsedTypeInfo{
		Request:    convertMessage(info.Request),
		Returns:    convertMessage(info.Returns),
		NewReturns: convertMessageFactory(info.NewReturns),
		NewRequest: convertMessageFactory(info.NewRequest),
	}
}
