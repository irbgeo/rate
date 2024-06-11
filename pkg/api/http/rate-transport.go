package http

import (
	"net/http"

	"github.com/irbgeo/rate/internal/controller"
)

type rateServerTransport struct {
}

func (s *rateServerTransport) decodeRequest(r *http.Request) (*controller.RateRequest, error) {
	req := &controller.RateRequest{}

	req.TokenIn = r.URL.Query().Get("tokenIn")
	if req.TokenIn == "" {
		return nil, errTokenInIsRequired
	}

	req.TokenOut = r.URL.Query().Get("tokenOut")
	if req.TokenOut == "" {
		return nil, errTokenOutIsRequired
	}

	return req, nil
}

func (s *rateServerTransport) encodeResponse(w http.ResponseWriter, r controller.RateResponse) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.Rate.String())) // nolint: errcheck
}
