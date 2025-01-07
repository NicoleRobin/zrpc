package handler

import (
	"context"
	pb "echo/api/echo"
	"echo/services/greeter"
)

type GreeterHandler struct {
}

func (GreeterHandler) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloRes, error) {
	return greeter.SayHello(ctx, req)
}
