PROTO_DIR = api/echo
VERSION := $(shell git describe --tags --dirty 2>/dev/null)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
TIME := $(shell date +"%Y-%m-%d.%H:%M:%S")
SETTINGS := github.com/nicolerobin/zrpc/settings
proto_set := $(wildcard $(PROTO_DIR)/*.proto)
FLAGS = -ldflags "-X $(SETTINGS).BuildVersion=$(VERSION) -X $(SETTINGS).BuildCommit=$(COMMIT) -X $(SETTINGS).BuildDate=$(TIME)"
pb_set := $(patsubst %.proto, %.pb.go, $(proto_set))

all: 
	echo

app:
	go build $(FLAGS)

proto: $(pb_set)
	@echo "Done"

$(pb_set): $(proto_set)
	protoc -I $(PROTO_DIR) --zrpc_out=plugins=zrpc+grpc:$(PROTO_DIR) $^

clean:
	rm -rf bin

.PHONY: all clean
