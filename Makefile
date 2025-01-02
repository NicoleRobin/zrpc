.PHONY: all clean

all: 
	protoc-gen-zrpc

protoc-gen-zrpc:
	go build -o cmds/protoc-gen-zrpc/protoc-gen-zrpc ./cmds/protoc-gen-zrpc

clean:
	rm -rf bin