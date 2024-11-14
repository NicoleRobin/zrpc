# protoc插件

## 用法
```shell
protoc -I example --zrpc_out=example example/example.proto
```

# 遇到的问题
1、`--zrpc_out: example.pb.go: generated file does not match prefix "example"`
原因：
居然是因为--zrpc_out指定的值必须是后续proto文件的前缀，不知道这是什么原因；

2、

# 参考文档
1、https://pkg.go.dev/google.golang.org/protobuf/compiler/protogen