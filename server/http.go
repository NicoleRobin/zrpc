package server

import (
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type httpHandler struct {
	s *grpc.Server
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}

func httpServer(l net.Listener, s *grpc.Server) error {
	hs := http.Server{
		Handler: &httpHandler{s},
	}

	return hs.Serve(l)
}
