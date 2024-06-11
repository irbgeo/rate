package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/irbgeo/rate/internal/controller"
)

type server struct {
	srv *http.Server
}

type svc interface {
	Rate(ctx context.Context, r controller.RateRequest) (controller.RateResponse, error)
}

func NewServer(svc svc) *server {
	router := mux.NewRouter()

	rateHandle(router, svc)

	return &server{
		srv: &http.Server{
			Handler: router,
		},
	}
}

func (s *server) ListenAndServe(port string) error {
	s.srv.Addr = ":" + port

	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
