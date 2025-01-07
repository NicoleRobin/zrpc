package rpc

import (
	"fmt"
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

func convertMessage(msg interface{}) proto.Message {
	switch v := msg.(type) {
	case proto.Message:
		return v
	default:
		return protoimpl.X.ProtoMessageV2Of(v)
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
		NewReturns: convertMessage(info.NewReturns),
		NewRequest: convertMessage(info.NewRequest),
	}
}
