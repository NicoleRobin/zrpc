package greeter

import (
	"context"

	"github.com/nicolerobin/zrpc/log"
	"go.uber.org/zap"

	pb "echo/api/echo"
)

func SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloRes, error) {
	log.Info(ctx, "entrance", zap.Any("req", req))
	return &pb.HelloRes{
		Message: "hello " + req.Name,
	}, nil
}
