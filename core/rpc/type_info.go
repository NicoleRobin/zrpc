package rpc

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

func RegisterTypeInfo(method string, info TypeInfo) {
}
