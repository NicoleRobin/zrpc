package greeter

import (
	"context"
	pb "echo/api/echo"
)

func SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloRes, error) {
	return &pb.HelloRes{
		Message: "hello " + req.Name,
	}, nil
}
