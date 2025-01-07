package server

import (
	"context"
	"fmt"
	"net"

	"github.com/nicolerobin/zrpc/config"
	"github.com/nicolerobin/zrpc/log"
)

// Start run server
func Start(ctx context.Context) error {
	addr := config.GetAddr()
	log.Infof(ctx, "starting server on %s", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf(ctx, "failed to listen: %v", err)
		return fmt.Errorf("failed to listen: %w", err)
	}
	s := newServer()
	return s.Serve(l)
}
